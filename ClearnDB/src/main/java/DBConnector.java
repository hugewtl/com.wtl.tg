import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

public class DBConnector {
        /**
         * Generate Connection
         */
        public Connection dbConnector(String dbtype, DBInit dbinit, LoggerLog log)
                        throws ClassNotFoundException, SQLException {
                // DBInit dbinit = new DBInit();
                // 1.regedit driver
                Class.forName(dbinit.getDbDriver());
                // 2.get connection
                if (dbtype.equals("oracle")) {
                        dbinit.connection = DriverManager.getConnection(
                                        "jdbc:oracle:thin:@//" + dbinit.getDbIp() + ":" + dbinit.getDbPort() + "/"
                                                        + dbinit.getDbSchema(),
                                        dbinit.getDbUser(), dbinit.getDbPasswd());
                }
                if (dbtype.equals("dm")) {
                        dbinit.connection = DriverManager.getConnection(
                                        "jdbc:dm://" + dbinit.getDbIp() + ":" + dbinit.getDbPort() + "?schema="
                                                        + dbinit.getDbSchema()
                                                        + "&compatibleMode=mysql&zeroDateTimeBehavior=convertToNull&useUnicode=true&characterEncoding=utf-8",
                                        dbinit.getDbUser(), dbinit.getDbPasswd());
                }

                if (dbtype.equals("mysql")) {
                        dbinit.connection = DriverManager.getConnection(
                                        "jdbc:mysql://" + dbinit.getDbIp() + ":" + dbinit.getDbPort()
                                                        + "/" + dbinit.getDbSchema()
                                                        + "?useSSL=false&characterEncoding=utf-8",
                                        dbinit.getDbUser(),
                                        dbinit.getDbPasswd());
                }
                log.logging(dbinit.connection + " The target database session inited and connected !");
                return dbinit.connection;
        }

        /* ConnectDB from the init connection with sqlstring for delete dataset */
        public boolean DelDT(Connection dbconn, String sqlString, String[] paramStr, LoggerLog log)
                        throws ClassNotFoundException, SQLException {
                /**
                 * predefine SQL and assignment,if oracle ,no ";" with ending!
                 */
                // log.logging(dbconn.toString().substring(0, 6));
                if (dbconn.toString().substring(0, 6).equals("oracle")) {
                        sqlString = sqlString.replace(";", "");
                }
                log.logging(dbconn + " The target database session connected and executing Delete SQL!");
                PreparedStatement preparedStatement = dbconn.prepareStatement(sqlString);
                // String[] paramStr = { "d01c508e95ec4a0b8ff65e53aa1f23a8", "GDL%" };
                if (paramStr.length != 0) {
                        this.setSqlStrParams(preparedStatement, paramStr, log);
                }
                /**
                 * 执行删除
                 */
                log.logging(sqlString);
                /**
                 * 更新数据要用executeUpdate()
                 */
                int delCount = preparedStatement.executeUpdate();
                log.logging(delCount + "");
                if (delCount > 0) {
                        log.logging("删除数据记录： " + delCount + " 条！");
                        return true;
                }
                return false;
        }

        /* ConnectDB from the init connection with sqlstring for getting dataset */
        public ResultSet QueryDB(Connection dbconn, String sqlString, String[] paramStr, LoggerLog log)
                        throws ClassNotFoundException, SQLException {
                /**
                 * predefine SQL and assignment,if oracle ,no ";" with ending!
                 */
                // log.logging(dbconn.toString().substring(0, 6));
                if (dbconn.toString().substring(0, 6).equals("oracle")) {
                        sqlString = sqlString.replace(";", "");
                }
                log.logging(dbconn + " The target database session connected and executing Query SQL!");
                PreparedStatement preparedStatement = dbconn.prepareStatement(sqlString);
                // String[] paramStr = { "d01c508e95ec4a0b8ff65e53aa1f23a8", "GDL%" };
                if (paramStr.length != 0) {
                        this.setSqlStrParams(preparedStatement, paramStr, log);
                }
                ResultSet rs = preparedStatement.executeQuery();
                log.logging(sqlString);
                return rs;
        }

        /* set String type parameters for the dynamic SQL */
        public void setSqlStrParams(PreparedStatement preparedStatement, String[] params, LoggerLog log) {
                int i = 0;
                for (i = 0; i < params.length; i++) {
                        try {
                                preparedStatement.setString(i + 1, params[i]);
                        } catch (SQLException e) {
                                log.logging(e.toString());
                        }
                }
        }

        /* set int type parameters for the dynamic SQL */
        public void setSqlIntParams(PreparedStatement preparedStatement, int[] params, LoggerLog log) {
                int i = 0;
                for (i = 0; i < params.length; i++) {
                        try {
                                preparedStatement.setInt(i + 1, params[i]);
                        } catch (SQLException e) {
                                log.logging(e.toString());
                        }
                }
        }
}
