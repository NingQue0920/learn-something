package main.demo.basic;

import sun.misc.Unsafe;

/**
 * 功能描述：
 *
 * @author liuchang49507
 * @date 2024/8/23
 */
public class CAS {
    public static void main(String[] args) {

        Unsafe.getUnsafe().compareAndSwapInt(1, 1, 1, 1);

    }
}
