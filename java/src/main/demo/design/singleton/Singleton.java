package main.demo.design.singleton;

public enum Singleton {
    INSTANCE;

    String name ;

    public void setName(String name) {
        this.name = name;
    }
    public void doSomething(){
        System.out.println("do something ");
    }
}
