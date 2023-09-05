package com.spring.tg;

// import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
@SpringBootApplication
public class TgApplication {

	public static void main(String[] args) {
		new TgApplication().logging();
		// SpringApplication.run(TgApplication.class, args);
	}

	// public class Logger {
    private static final Logger logger = LoggerFactory.getLogger(TgApplication.class.getName());
    public void logging() {
        // logger.debug("debug message");
        logger.info("info message"); 
        // logger.warn("warn message");
        // logger.error("error message");
        // logger.fatal("fatal message");
    }
// }
}
