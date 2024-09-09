package main.demo.basic;

public class ThreadLocalDemo {

    static class BigObject {
        byte[] data = new byte[1024 * 1024];
    }


    private static final ThreadLocal<BigObject> threadLocal =    new ThreadLocal<BigObject>();


    public static void main(String[] args) throws InterruptedException {




        Thread.sleep(10000);
        Thread thread = new Thread(() -> {
            threadLocal.set(new BigObject());
            System.out.println("Thread local value set ");
        });
        thread.start();

        threadLocal.remove();



    }
}
