import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.Properties;

public class ReadParams {
    /**
     * 将properties作为类对象，可直接访问
     */
    Properties properties = new Properties();

    public void readParams(LoggerLog log) {
        InputStream inputStream = this.getClass().getResourceAsStream("export.properties");
        try {
            /**
             * 解决中文乱码问题
             */
            BufferedReader bfReader = new BufferedReader(new InputStreamReader(inputStream, "UTF-8"));
            properties.load(bfReader);
        } catch (IOException e) {
            log.logging(e.toString());
        }
        // log.logging("==============================================");
        // String property = properties.getProperty("sqlfile");
        // log.logging("property = " + property);
        properties.list(System.out);
        // log.logging("==============================================");
    }
}
