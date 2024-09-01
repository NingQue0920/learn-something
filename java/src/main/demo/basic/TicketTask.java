package main.demo.basic;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class TicketTask implements Runnable{


    int ticket = 5 ;
    private static final Lock lock = new ReentrantLock();
    @Override
    public void run() {
        try {
            while (true){
                lock.lock();
                if (ticket<=0){
                    System.out.println(Thread.currentThread().getName() + " 售空 ---");
                    break;
                }else {
                    ticket--;
                    System.out.println(Thread.currentThread().getName() + " 卖出第" + (ticket +1)+ " 张");
                }
            }
        }finally {
            lock.unlock();
            lock.unlock();lock.unlock();lock.unlock();lock.unlock();lock.unlock();

            System.out.println(Thread.currentThread().getName() + " 释放锁");
        }
    }

    public static void main(String[] args) {
        TicketTask ticketTask = new TicketTask();
        Thread t1 = new Thread(ticketTask, "一号窗口");
        Thread t2 = new Thread(ticketTask, "二号窗口");
        t1.start();
        t2.start();
    }
}
