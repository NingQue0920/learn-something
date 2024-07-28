package main.demo.spring;


import org.springframework.beans.BeansException;
import org.springframework.beans.factory.*;
import org.springframework.beans.factory.config.BeanPostProcessor;
import org.springframework.beans.factory.xml.XmlBeanDefinitionReader;
import org.springframework.boot.autoconfigure.AutoConfigurationImportListener;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Import;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;

@EnableAutoConfiguration()
public class BeanLifeCycle {

    public static void main(String[] args) {
        AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(BeanLifeCycle.class);

        MyBean myBean = context.getBean(MyBean.class);
        myBean.doSomething();
        context.close();
    }

    @Component
    public static class MyBean extends XmlBeanDefinitionReader implements InitializingBean, DisposableBean, BeanNameAware, BeanFactoryAware, ApplicationContextAware, BeanPostProcessor {

        private String beanName;

        public MyBean() {
            super( new AnnotationConfigApplicationContext(  ));
            System.out.println("1. Bean实例化 - 构造方法被调用");
        }

        @Override
        public void setBeanName(String name) {
            this.beanName = name;
            System.out.println("3. BeanNameAware's setBeanName 被调用: " + name);
        }

        @Override
        public void setBeanFactory(BeanFactory beanFactory) throws BeansException {
            System.out.println("4. BeanFactoryAware's setBeanFactory 被调用");
        }

        @Override
        public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
            System.out.println("5. ApplicationContextAware's setApplicationContext 被调用");
        }

        @PostConstruct
        public void postConstruct() {
            System.out.println("6. @PostConstruct 注解方法被调用");
        }

        @Override
        public void afterPropertiesSet() throws Exception {
            System.out.println("7. InitializingBean's afterPropertiesSet 被调用");
        }

        public void customInit() {
            System.out.println("8. 自定义初始化方法 customInit 被调用");
        }

        public void doSomething() {
            System.out.println("9. Bean 正在使用中 - doSomething 方法被调用");
        }

        @PreDestroy
        public void preDestroy() {
            System.out.println("10. @PreDestroy 注解方法被调用");
        }

        @Override
        public void destroy() throws Exception {
            System.out.println("11. DisposableBean's destroy 被调用");
        }

        public void customDestroy() {
            System.out.println("12. 自定义销毁方法 customDestroy 被调用");
        }


    }

    @Bean(initMethod = "customInit", destroyMethod = "customDestroy")
    public MyBean myBean() {
        return new MyBean();
    }

    @Component
    public static class MyBeanPostProcessor implements BeanPostProcessor {
        @Override
        public Object postProcessBeforeInitialization(Object bean, String beanName) throws BeansException {
            if (bean instanceof MyBean) {
                System.out.println("6. BeanPostProcessor's postProcessBeforeInitialization 被调用");
            }
            return bean;
        }

        @Override
        public Object postProcessAfterInitialization(Object bean, String beanName) throws BeansException {
            if (bean instanceof MyBean) {
                System.out.println("8. BeanPostProcessor's postProcessAfterInitialization 被调用");
            }
            return bean;
        }
    }
}
