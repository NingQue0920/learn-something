package main.demo.design.strategy;



public class Main {
    public static void main(String[] args) {
        int a = 10 ;
        int b = 5;
        Operation.ADD.execute(a, b);
        Operation.SUBTRACT.execute(a, b);
        Operation.MULTIPLY.execute(a, b);
        Operation.DIVIDE.execute(a, b);
    }
}
