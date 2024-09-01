package main.demo.design.singleton;

import java.io.*;
import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;

public class Main {
    public static void main(String[] args) throws NoSuchMethodException, InvocationTargetException, InstantiationException, IllegalAccessException, IOException, ClassNotFoundException {
        Singleton.INSTANCE.doSomething();
        Singleton.INSTANCE.setName("name");
        System.out.println(Singleton.INSTANCE.name);

        // 防止反射
        Constructor<?> constructor = Singleton.class.getDeclaredConstructor();
        constructor.setAccessible(true);
        constructor.newInstance();

        // 防止序列化
        ObjectOutputStream out = new ObjectOutputStream(new FileOutputStream("singleton.ser"));
        out.writeObject(Singleton.INSTANCE);
        out.close();

        ObjectInputStream in = new ObjectInputStream(new FileInputStream("singleton.ser"));
        Singleton instance = (Singleton) in.readObject();
        in.close();

        System.out.println(instance == Singleton.INSTANCE); // Should print true
        instance.doSomething();


    }
}
