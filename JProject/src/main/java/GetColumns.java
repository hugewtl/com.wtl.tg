import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.LinkedHashMap;

import java.sql.ResultSetMetaData;

public class GetColumns {

    /**
     * get and store column name and data type from ResultSet
     */
    LinkedHashMap<String, String> fieldmap = new LinkedHashMap<>();

    public LinkedHashMap<String, String> getColumns(ResultSet resultSet, LoggerLog log) {
        ResultSetMetaData metaData = null;
        try {
            metaData = resultSet.getMetaData();
        } catch (SQLException e) {
            log.logging(e.toString());
        }
        int columnCount = 0;
        try {
            columnCount = metaData.getColumnCount();
            for (int i = 1; i <= columnCount; i++) {
                String columnName = metaData.getColumnName(i);
                String columnType = metaData.getColumnTypeName(i);
                fieldmap.put(columnName, columnType);
                // log.logging("Column Name: " + columnName + ", Column Type: " + columnType);
            }
        } catch (SQLException e) {
            log.logging(e.toString());
        }

        return fieldmap;
    }

    /**
     * contact fields data for SQL statement with NULL,commoa,single quotation mark
     */
    public String appendVals(String rec, String val) {
        /**
         * 将字段中包含"''替换成"''"
         */

        if (val.contains("'")) {
            val = val.replaceAll("'", "''");
        }
        // 将字段拼接成字符串，处理字段为空输出逻辑+字段拼接逻辑
        if (!val.equals("NULL")) {
            val = "'" + val + "'";
        }
        /**
         * 处理Oracle时间戳字段匹配时间戳格式
         */
        if (val.startsWith("'20") && val.endsWith(".000000'")) {
            String timestamp = "TIMESTAMP";
            val = timestamp + " " + val;
        }
        if (!rec.equals("NULL")) {
            rec = rec + "," + val;
        } else {
            rec = val;
        }
        return rec;
    }

}
