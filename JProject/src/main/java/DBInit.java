import java.sql.Connection;

public class DBInit {
    java.sql.Statement statement;
    Connection connection;
    /**
     * 数据库连接访问参数
     */
    String dbDriver;
    String dbIp;
    String dbUser;
    String dbPort;
    String dbPasswd;
    String dbSchema;

    /**
     * @return the dbSchema
     */
    public String getDbSchema() {
        return dbSchema;
    }

    /**
     * @param dbSchema the dbSchema to set
     */
    public void setDbSchema(String dbSchema) {
        this.dbSchema = dbSchema;
    }

    /**
     * @return the dbDriver
     */
    public String getDbDriver() {
        return dbDriver;
    }

    /**
     * @param dbDriver the dbDriver to set
     */
    public void setDbDriver(String dbDriver) {
        this.dbDriver = dbDriver;
    }

    /**
     * @return the statement
     */
    public java.sql.Statement getStatement() {
        return statement;
    }

    /**
     * @param statement the statement to set
     */
    public void setStatement(java.sql.Statement statement) {
        this.statement = statement;
    }

    /**
     * @return the connection
     */
    public Connection getConnection() {
        return connection;
    }

    /**
     * @param connection the connection to set
     */
    public void setConnection(Connection connection) {
        this.connection = connection;
    }

    /**
     * @return the dbIp
     */
    public String getDbIp() {
        return dbIp;
    }

    /**
     * @param dbIp the dbIp to set
     */
    public void setDbIp(String dbIp) {
        this.dbIp = dbIp;
    }

    /**
     * @return the dbUser
     */
    public String getDbUser() {
        return dbUser;
    }

    /**
     * @param dbUser the dbUser to set
     */
    public void setDbUser(String dbUser) {
        this.dbUser = dbUser;
    }

    /**
     * @return the dbPort
     */
    public String getDbPort() {
        return dbPort;
    }

    /**
     * @param dbPort the dbPort to set
     */
    public void setDbPort(String dbPort) {
        this.dbPort = dbPort;
    }

    /**
     * @return the dbPasswd
     */
    public String getDbPasswd() {
        return dbPasswd;
    }

    /**
     * @param dbPasswd the dbPasswd to set
     */
    public void setDbPasswd(String dbPasswd) {
        this.dbPasswd = dbPasswd;
    }
}
