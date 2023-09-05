import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;

public class GetSQLFile {
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
         * Get Resultset from the target SQL
         */
        GetTargetSet gts = new GetTargetSet();

        /**
         * Export the target rules
         */
        String[] paramStr1 = new String[1];
        /* Get ID from AM_RULE_SET with set_name and states=1 */
        String sqlString_pre1 = "SELECT ID FROM AM_RULE_SET WHERE SET_NAME = ? AND STATES = 1;";
        paramStr1[0] = rp.properties.getProperty("set_name");
        try {
            Connection conn = dbconnector.dbConnector(rp.properties.getProperty("dbtype"), dbi, log);
            /* Get ID from am_ruleset_history with set_id and hist_verson */
            ResultSet rs1 = dbconnector.QueryDB(conn, sqlString_pre1, paramStr1, log);
            String sqlString_pre2 = "SELECT ID FROM AM_RULESET_HISTORY WHERE SET_ID = ? AND HIST_VERSION = ?;";
            String[] paramStr2 = new String[2];
            while (rs1.next()) {
                paramStr2[0] = rs1.getString(1);
                log.logging("AM_RULESET_HISTORY获取目标SET_ID：" + paramStr2[0]);
            }
            rs1.close();
            /*
             * Get ResultSet from am_rule with RULESET_HISTORY_ID and RULE_NO prefix label
             */
            gts.tableName = "AM_RULE";
            paramStr2[1] = rp.properties.getProperty("set_hist_version");
            ResultSet rs2 = dbconnector.QueryDB(conn, sqlString_pre2, paramStr2, log);
            String sqlString = "SELECT * FROM AM_RULE WHERE RULESET_HISTORY_ID = ? AND RULE_NO LIKE ?;";
            String[] paramStr = new String[2];
            while (rs2.next()) {
                paramStr[0] = rs2.getString(1);
                log.logging("AM_RULE获取目标RULESET_HISTORY_ID：" + paramStr[0]);
            }
            rs2.close();
            paramStr[1] = rp.properties.getProperty("rule_no_prefix");

            if (rp.properties.get("if_all_rule").equals("yes")) {
                /*
                 * 全场景规则导出
                 */
                String sqlString_ar = "SELECT * FROM AM_RULE WHERE RULESET_HISTORY_ID = ?;";
                String[] paramStr_ar = new String[1];
                paramStr_ar[0] = paramStr[0];
                if (gts.getQuerySet(conn, sqlString_ar, paramStr_ar, dbconnector, log)) {
                    log.logging("<RULESET_HISTORY_ID:" + paramStr_ar[0] + ">"
                            + " 查询结果集导出成功！");
                }
            } else {
                if (gts.getQuerySet(conn, sqlString, paramStr, dbconnector, log)) {
                    log.logging("<RULESET_HISTORY_ID:" + paramStr[0] + " ; rule_no_prefix:" + paramStr[1] + ">"
                            + " 查询结果集导出成功！");
                }
            }
            /**
             * am_rule中提取到指标，名单、标签、字段数据
             */
            log.logging(gts.tableName + " 中提取" + getLen(MyAmRule.calId) + "个指标ID:" + MyAmRule.calId);
            log.logging(gts.tableName + " 中提取" + getLen(MyAmRule.fieldId) + "个业务字段:" + MyAmRule.fieldId);
            log.logging(gts.tableName + " 中提取" + getLen(MyAmRule.enumId) + "个枚举字段:" + MyAmRule.enumId);
            log.logging(gts.tableName + " 中提取" + getLen(MyAmRule.nameListId) + "个名单ID:" +
                    MyAmRule.nameListId);
            log.logging(gts.tableName + " 中提取" + getLen(GetTargetSet.labelId) + "个标签ID:" +
                    GetTargetSet.labelId);
            /**
             * 从am_rule中关联到的主体数据再次提取：指标、名单、标签、字段
             */
            // 提取指标,首先对条件进行判空处理
            if (MyAmRule.calId.isEmpty()) {
                log.logging("未获取到统计的指标ID数据");
            } else {
                /**
                 * 依据查询出的指标ID和字段ID与上次是否完全一致，判断是否循环查出所有结果
                 */
                boolean isOver = false;
                for (; !isOver;) {
                    String calId = MyAmRule.calId;
                    String fieldId = MyAmRule.fieldId;
                    String enumId = MyAmRule.enumId;
                    if (!MyAmRule.calId.isEmpty()) {
                        gts.tableName = "AM_INDICATORS";
                        String sqlString_calId = "SELECT * FROM " + gts.tableName + " WHERE ID IN (" + MyAmRule.calId
                                + ") AND 1 = ?;";
                        /**
                         * 为了复用函数，传个恒真值给占位符
                         */
                        String[] paramStr_calId = new String[1];
                        paramStr_calId[0] = "1";
                        if (gts.getQuerySet(conn, sqlString_calId, paramStr_calId, dbconnector, log)) {
                            /**
                             * am_indicators中提取到指标、字段数据
                             */

                            log.logging(gts.tableName + " 中新增提取" + (getLen(MyAmRule.calId) - getLen(calId)) + "个指标ID:"
                                    + MyAmRule.calId);
                            log.logging(
                                    gts.tableName + " 中新增提取" + (getLen(MyAmRule.fieldId) - getLen(fieldId)) + "个业务字段:"
                                            + MyAmRule.fieldId);
                            log.logging(gts.tableName + " 中新增提取" + (getLen(MyAmRule.enumId) - getLen(enumId)) + "个枚举字段:"
                                    + MyAmRule.enumId);
                            log.logging(gts.tableName + " 提取< " + getLen(MyAmRule.calId) + " 个指标ID:" + MyAmRule.calId
                                    + ">查询结果集导出成功！");
                        }
                    }

                    if (!MyAmRule.fieldId.isEmpty()) {
                        // 提取业务字段 和 枚举字段
                        gts.tableName = "AM_ENUMERATE_FIELD";
                        if (!MyAmRule.enumId.isEmpty()) {
                            MyAmRule.fieldId = MyAmRule.fieldId + "," + MyAmRule.enumId;
                        }
                        String sqlString_fieldId = "SELECT * FROM " + gts.tableName + " WHERE FIELD_VALUE IN ("
                                + MyAmRule.fieldId
                                + ") AND 1 = ?;";
                        /**
                         * 为了复用函数，传个恒真值给占位符
                         */
                        String[] paramStr_fieldId = new String[1];
                        paramStr_fieldId[0] = "1";
                        if (gts.getQuerySet(conn, sqlString_fieldId, paramStr_fieldId, dbconnector, log)) {
                            log.logging(gts.tableName + " 提取< " + getLen(MyAmRule.fieldId) + " 个字段值:" + MyAmRule.fieldId
                                    + ">查询结果集导出成功！");
                        }
                    } else {
                        log.logging("未获取到业务字段或枚举字段！");
                    }
                    /**
                     * 判断结束依据,分析多次统计的必要性
                     */
                    if (calId.equals(MyAmRule.calId) && fieldId.equals(MyAmRule.fieldId)) {
                        isOver = true;
                    }
                }
            }

            if (!MyAmRule.nameListId.isEmpty()) {
                // 提取名单
                gts.tableName = "AM_LIST_NAME";
                String sqlString_nameListId = "SELECT * FROM " + gts.tableName + " WHERE ID IN (" + MyAmRule.nameListId
                        + ") AND 1 = ?;";
                /**
                 * 为了复用函数，传个恒真值给占位符
                 */
                String[] paramStr_nameListId = new String[1];
                paramStr_nameListId[0] = "1";
                if (gts.getQuerySet(conn, sqlString_nameListId, paramStr_nameListId, dbconnector, log)) {
                    log.logging(gts.tableName + " 提取< " + getLen(MyAmRule.nameListId) + " 个名单集ID:" + MyAmRule.nameListId
                            + ">查询结果集导出成功！");
                }
            } else {
                log.logging("未获取到名单ID数据！");
            }

            if (!GetTargetSet.labelId.isEmpty()) {
                // 提取标签
                gts.tableName = "AM_SYS_PARAMS";
                String sqlString_labelId = "SELECT * FROM " + gts.tableName + " WHERE ID IN (" + GetTargetSet.labelId
                        + ") AND 1 = ?;";
                /**
                 * 为了复用函数，传个恒真值给占位符
                 */
                String[] paramStr_labelId = new String[1];
                paramStr_labelId[0] = "1";
                if (gts.getQuerySet(conn, sqlString_labelId, paramStr_labelId, dbconnector, log)) {
                    log.logging(
                            gts.tableName + " 提取< " + getLen(GetTargetSet.labelId) + " 个标签ID:" + GetTargetSet.labelId
                                    + ">查询结果集导出成功！");
                }
            } else {
                log.logging("未获取到标签ID数据！");
            }

            /**
             * 导出新增规则模板和参数
             */
            if (!rp.properties.getProperty("indic_template_name").isEmpty()) {
                gts.tableName = "AM_INDIC_TMPL";
                String sqlString_tmpl = "SELECT * FROM AM_INDIC_TMPL WHERE temp_name IN ("
                        + rp.properties.getProperty("indic_template_name") + ") AND 1 = ?;";
                String[] paramStr_tmpl = new String[1];
                paramStr_tmpl[0] = "1";
                if (gts.getQuerySet(conn, sqlString_tmpl, paramStr_tmpl, dbconnector, log)) {
                    log.logging(gts.tableName + " 提取<模板名称:" + rp.properties.getProperty("indic_template_name")
                            + ">查询结果集导出成功！");
                }

                /*
                 * 导出指标模板参数
                 */
                if (!GetTargetSet.indic_tmpl_ids.isEmpty()) {
                    gts.tableName = "AM_INDIC_TMPL_PARAM";
                    String sqlString_tmpl_param = "SELECT * FROM AM_INDIC_TMPL_PARAM WHERE TEMPLATE_ID IN ("
                            + GetTargetSet.indic_tmpl_ids + ") AND 1 = ?;";
                    String[] paramStr_tmpl_param = new String[1];
                    paramStr_tmpl_param[0] = "1";
                    if (gts.getQuerySet(conn, sqlString_tmpl_param, paramStr_tmpl_param, dbconnector, log)) {
                        log.logging(gts.tableName + " 提取<指标模板参数:" + rp.properties.getProperty("indic_template_name")
                                + ">查询结果集导出成功！");
                    }
                }

            }

            /**
             * 导出新增方法和参数
             */
            if (!rp.properties.getProperty("indic_method_name").isEmpty()) {
                gts.tableName = "AM_METHOD";
                String sqlString_mth = "SELECT * FROM AM_METHOD WHERE MTHOD_NAME IN ("
                        + rp.properties.getProperty("indic_method_name") + ")  AND 1 = ?;";
                String[] paramStr_mth = new String[1];
                paramStr_mth[0] = "1";
                if (gts.getQuerySet(conn, sqlString_mth, paramStr_mth, dbconnector, log)) {
                    log.logging(gts.tableName + " 提取<方法名称:" + rp.properties.getProperty("indic_method_name")
                            + ">查询结果集导出成功！");
                }

                /*
                 * 导出方法模板参数
                 */
                if (!GetTargetSet.mth_tmpl_ids.isEmpty()) {
                    gts.tableName = "AM_METHOD_PARAM";
                    String sqlString_mth_param = "SELECT * FROM AM_METHOD_PARAM WHERE MID IN ("
                            + GetTargetSet.mth_tmpl_ids + ") AND 1 = ?;";
                    String[] paramStr_mth_param = new String[1];
                    paramStr_mth_param[0] = "1";
                    if (gts.getQuerySet(conn, sqlString_mth_param, paramStr_mth_param, dbconnector, log)) {
                        log.logging(gts.tableName + " 提取<指标模板参数:" + rp.properties.getProperty("indic_method_name")
                                + ">查询结果集导出成功！");
                    }
                }

            }
            /*
             * 关闭数据库连接
             */
            conn.close();
        } catch (SQLException e) {
            log.logging(e.toString());
        } catch (ClassNotFoundException e) {
            log.logging(e.toString());
        }

    }

    public static int getLen(String valId) {
        if (valId.length() == 0) {
            return 0;
        }
        return valId.split(",").length;
    }

    /**
     * 字符串数组转换为set用于集合计算
     */
    // public static HashSet<String> getSetFromString(String srcStr) {
    // if (!srcStr.isEmpty()) {
    // String[] strArr = srcStr.split(",");
    // ArrayList<String> list = (ArrayList<String>) Arrays.asList(strArr);
    // HashSet<String> set = new HashSet<>(list);
    // return set;
    // }
    // return null;
    // }
}
