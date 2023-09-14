import java.sql.Connection;
import java.sql.SQLException;

public class CleanDB {
    public static String tableName = "";
    public static String colName = "";

    public static void main(String[] args) {
        /**
         * Init logger and the application configurations
         */
        LoggerLog log = new LoggerLog();
        ReadParams rp = new ReadParams();
        DBConnector dbconnector = new DBConnector();
        DBInit dbi = new DBInit();
        /**
         * Load propertiy here will get value directly
         */
        rp.readParams(log);
        /**
         * Init DB connection with getting necessary parameters
         */
        dbi.setDbDriver(rp.properties.getProperty("dbdriver"));
        dbi.setDbIp(rp.properties.getProperty("ip"));
        dbi.setDbPort(rp.properties.getProperty("port"));
        dbi.setDbUser(rp.properties.getProperty("dbuser"));
        dbi.setDbPasswd(rp.properties.getProperty("dbpasswd"));
        dbi.setDbSchema(rp.properties.getProperty("schema"));
        /**
         * 定义要清理的数据表，删除条件
         */
        String[] paramStr = new String[2];
        tableName = rp.properties.getProperty("tableNameParam");
        colName = rp.properties.getProperty("colNameParam");
        paramStr[0] = "1";
        paramStr[1] = rp.properties.getProperty("colNameCondition");
        String sqlString = "DELETE FROM " + tableName + " WHERE 1 = ? " + " AND " + colName + " = ?;";
        try {
            Connection conn = dbconnector.dbConnector(rp.properties.getProperty("dbtype"), dbi, log);
            if (DeleteDt(conn, sqlString, paramStr, dbconnector, log)) {
                log.logging(tableName + " 表数据删除目标数据完成 !");
            } else {
                log.logging(tableName + " 未执行删除 !");
            }
        } catch (SQLException e) {
            log.logging(e.toString());
        } catch (ClassNotFoundException e) {
            log.logging(e.toString());
        }

    }

    public static boolean DeleteDt(Connection conn, String sqlString, String[] paramStr, DBConnector dbconnector,
            LoggerLog log) {
        /* 返回执行状态，查询am_rule 导出sql */
        // boolean isOracle = conn.toString().substring(0, 6).equals("oracle");
        // boolean isMySQL = conn.toString().substring(0, 9).equals("com.mysql");
        // boolean isDM = conn.toString().substring(0, 2).equals("dm");
        try {
            return dbconnector.DelDT(conn, sqlString, paramStr, log);
        } catch (SQLException e) {
            log.logging(e.toString());
        } catch (ClassNotFoundException e) {
            log.logging(e.toString());
        }
        return false;

    }
}
