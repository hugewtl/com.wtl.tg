import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class LoggerLog {
    private static final Logger logger = LoggerFactory.getLogger(LoggerLog.class);

    public void logging(String loginfo) {
        logger.info(loginfo);
        // logger.debug(loginfo);
    }
}
