package main.demo.design.strategy;

public enum Operation implements Strategy{
    ADD {
        @Override
        public void execute(int a, int b) {
            System.out.println(a + b);
        }
    },
    SUBTRACT {
        @Override
        public void execute(int a, int b) {
            System.out.println(a - b);
        }
    },
    MULTIPLY {
        @Override
        public void execute(int a, int b) {
            System.out.println(a * b);
        }
    },
    DIVIDE {
        @Override
        public void execute(int a, int b) {
            System.out.println(a / b);
        }
    }
}