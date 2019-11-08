package misc

import (
	"github.com/zxh0/jvm.go/rtda"
	"github.com/zxh0/jvm.go/rtda/heap"
	"github.com/zxh0/jvm.go/vmutils"
)

func init() {
	_unsafe(allocateInstance, "allocateInstance", "(Ljava/lang/Class;)Ljava/lang/Object;")
	_unsafe(defineClass, "defineClass", "(Ljava/lang/String;[BIILjava/lang/ClassLoader;Ljava/security/ProtectionDomain;)Ljava/lang/Class;")
	_unsafe(shouldBeInitialized0, "shouldBeInitialized0", "(Ljava/lang/Class;)Z")
	_unsafe(ensureClassInitialized0, "ensureClassInitialized0", "(Ljava/lang/Class;)V")
	_unsafe(staticFieldOffset0, "staticFieldOffset0", "(Ljava/lang/reflect/Field;)J")
	_unsafe(staticFieldBase0, "staticFieldBase0", "(Ljava/lang/reflect/Field;)Ljava/lang/Object;")
}

// public native Object allocateInstance(Class<?> type) throws InstantiationException;
// (Ljava/lang/Class;)Ljava/lang/Object;
func allocateInstance(frame *rtda.Frame) {
	classObj := frame.GetRefVar(1)

	class := classObj.GetGoClass()
	obj := class.NewObj()

	frame.PushRef(obj)
}

// public native Class defineClass(String name, byte[] b, int off, int len,
//  		ClassLoader loader, ProtectionDomain protectionDomain)
// (Ljava/lang/String;[BIILjava/lang/ClassLoader;Ljava/security/ProtectionDomain;)Ljava/lang/Class;
func defineClass(frame *rtda.Frame) {
	nameObj := frame.GetRefVar(1)
	byteArr := frame.GetRefVar(2)
	off := frame.GetIntVar(3)
	_len := frame.GetIntVar(4)
	//loaderObj := frame.GetRefVar(5)
	//pd := frame.GetRefVar(6)

	name := nameObj.JSToGoStr()
	name = vmutils.DotToSlash(name)
	data := byteArr.GetGoBytes()
	data = data[off : off+_len]

	// todo
	class := frame.GetClassLoader().DefineClass(name, data)
	frame.PushRef(class.JClass)
}

// public native boolean shouldBeInitialized0(Class<?> c);
// (Ljava/lang/Class;)V
func shouldBeInitialized0(frame *rtda.Frame) {
	// this := frame.GetRefVar(0)
	classObj := frame.GetRefVar(1)

	goClass := classObj.GetGoClass()
	ret := goClass.InitializationNotStarted() // TODO
	frame.PushBoolean(ret)
}

// public native void ensureClassInitialized0(Class<?> c);
// (Ljava/lang/Class;)V
func ensureClassInitialized0(frame *rtda.Frame) {
	// this := frame.GetRefVar(0)
	classObj := frame.GetRefVar(1)

	goClass := classObj.GetGoClass()
	if goClass.InitializationNotStarted() {
		// undo ensureClassInitialized0()
		frame.RevertNextPC()
		// init
		frame.Thread.InitClass(goClass)
	}
}

// public native long staticFieldOffset0(Field f);
// (Ljava/lang/reflect/Field;)J
func staticFieldOffset0(frame *rtda.Frame) {
	// frame.GetRefVar(0) // this
	fieldObj := frame.GetRefVar(1)

	offset := fieldObj.GetFieldValue("slot", "I").IntValue()
	frame.PushLong(int64(offset))
}

// public native Object staticFieldBase0(Field f);
// (Ljava/lang/reflect/Field;)Ljava/lang/Object;
func staticFieldBase0(frame *rtda.Frame) {
	// frame.GetRefVar(0) // this
	fieldObj := frame.GetRefVar(1)

	goField := _getGoField(fieldObj)
	obj := goField.Class.AsObj()

	frame.PushRef(obj)
}

func _getGoField(fieldObj *heap.Object) *heap.Field {
	extra := fieldObj.Extra
	if extra != nil {
		return extra.(*heap.Field)
	}

	root := fieldObj.GetFieldValue("root", "Ljava/lang/reflect/Field;").Ref
	return root.Extra.(*heap.Field)
}
