import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.Reader;
import java.io.UnsupportedEncodingException;
import java.sql.Clob;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import com.google.gson.JsonSyntaxException;

public class GetTargetSet {
    String tableName;
    /*
     * 提取风险标签存储
     */
    public static String labelId = "";
    /**
     * 提取指标、方法模板ID值
     */
    public static String indic_tmpl_ids = "";
    public static String mth_tmpl_ids = "";

    /**
     * @return the tableName
     */
    public String getTableName() {
        return tableName;
    }

    /**
     * @param tableName the tableName to set
     */
    public void setTableName(String tableName) {
        this.tableName = tableName;
    }

    /* 返回执行状态，查询am_rule 导出sql */
    public boolean getQuerySet(Connection conn, String sqlString, String[] paramStr, DBConnector dbconnector,
            LoggerLog log) {
        boolean isOracle = conn.toString().substring(0, 6).equals("oracle");
        boolean isMySQL = conn.toString().substring(0, 9).equals("com.mysql");
        boolean isDM = conn.toString().substring(0, 2).equals("dm");
        try {
            ResultSet rs = dbconnector.QueryDB(conn, sqlString, paramStr, log);
            /**
             * if sqlfile directory is not exists ,create it
             */
            File sqlfile = new File("sqlfile");
            if (!sqlfile.exists()) {
                if (sqlfile.mkdir()) {
                    log.logging(sqlfile + " 目录创建成功");
                }
            }
            PrintWriter fos = new PrintWriter(sqlfile + "/" + this.tableName + ".sql", "utf-8");
            /**
             * generate the fiedlds fields sorted map
             */
            GetColumns gcl = new GetColumns();
            gcl.getColumns(rs, log);
            /**
             * generate the fiedlds columns
             */
            String cloStr = "";
            for (String columnName : gcl.fieldmap.keySet()) {
                if (cloStr.length() == 0) {
                    cloStr = columnName.toUpperCase();
                } else {
                    cloStr = cloStr + "," + columnName.toUpperCase();
                }
            }
            // log.logging(cloStr);
            /**
             * get column data from ResultSet
             */
            StringBuilder sqlBuilder = new StringBuilder();
            String dataStr = "";
            while (rs.next()) {
                if (isDM) {
                    dataStr = dealColDtForDM(rs, log, gcl);
                }
                if (isOracle) {
                    dataStr = dealColDtForOracle(rs, log, gcl);
                }
                if (isMySQL) {
                    dataStr = dealColDtForMySQL(rs, log, gcl);
                }
                sqlBuilder = new StringBuilder("INSERT INTO ")
                        .append(this.tableName + " (")
                        .append(cloStr).append(") VALUES (")
                        .append(dataStr)
                        .append(");");
                fos.println(sqlBuilder.toString());
                sqlBuilder = null;
            }
            fos.close();
            rs.close();

        } catch (SQLException e) {
            log.logging(e.toString());
            return false;
        } catch (ClassNotFoundException e) {
            log.logging(e.toString());
            return false;
        } catch (FileNotFoundException e) {
            log.logging(e.toString());
            return false;
        } catch (UnsupportedEncodingException e) {
            log.logging(e.toString());
            return false;
        }
        return true;
    }

    /*
     * 处理dm字段类型数据
     */
    public String dealColDtForDM(ResultSet rs, LoggerLog log, GetColumns gcl) throws SQLException {
        /**
         * define NULL for String constant
         */
        String dataStr = "NULL";
        for (String columnName : gcl.fieldmap.keySet()) {
            String dataType = gcl.fieldmap.get(columnName);
            // 处理DM的数据类型
            /**
             * 对TIMESTAMP类型处理
             */
            if (dataType.equals("TIMESTAMP")) {
                if (rs.getTimestamp(columnName, null) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr, rs.getTimestamp(columnName, null).toString());
                    // log.logging(dataStr);
                }
            }
            /**
             * 对CLOB类型处理--DmdbNClob
             */
            // log.logging(dataType);
            if (dataType.equals("CLOB") || dataType.equals("TEXT")) {
                if (rs.getClob(columnName) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    try {
                        Clob clob = rs.getClob(columnName);
                        String clobString = this.ClobToString(clob);
                        dataStr = gcl.appendVals(dataStr, clobString);
                        if (this.tableName.equals("AM_RULE")) {
                            /* 遍历解析rule_josn */
                            if (columnName.equals("RULE_JSON")) {
                                MyAmRule.parseRuleJson(clobString, log);
                            }
                        }
                        /*
                         * 嵌套指标提取:指标表
                         */
                        if (this.tableName.equals("AM_INDICATORS")) {
                            if (columnName.equals("INDICATORS_PARAM")) {
                                MyIndicator.parseIndicParam(clobString, log);
                            }
                        }
                    } catch (JsonSyntaxException e) {
                        log.logging(e.toString());
                    } catch (SQLException e) {
                        log.logging(e.toString());
                    } catch (IOException e) {
                        log.logging(e.toString());
                    }
                }
            }
            /**
             * 对String类型处理
             */
            if (dataType.contains("CHAR")) {
                if (rs.getString(columnName) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr, rs.getString(columnName));
                    /* 提取风险标签 */
                    if (columnName.equals("LABEL_ID")) {
                        if (!rs.getString(columnName).isEmpty()) {
                            if (labelId.isEmpty()) {
                                labelId = "'" + rs.getString(columnName) + "'";
                            } else {
                                labelId = MyAmRule.appendValsSingle(labelId, rs.getString(columnName));
                            }
                        }
                    }

                    /*
                     * 提取指标模板ID
                     */
                    if (this.tableName.equals("AM_INDIC_TMPL")) {
                        if (columnName.equals("ID")) {
                            if (indic_tmpl_ids.isEmpty() || indic_tmpl_ids == null) {
                                indic_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                indic_tmpl_ids = MyAmRule.appendValsSingle(indic_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }

                    /*
                     * 提取方法模板ID
                     */
                    if (this.tableName.equals("AM_METHOD")) {
                        if (columnName.equals("ID")) {
                            if (mth_tmpl_ids.isEmpty() || mth_tmpl_ids == null) {
                                mth_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                mth_tmpl_ids = MyAmRule.appendValsSingle(mth_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }

                }
            }
        }
        return dataStr;
    }

    /*
     * 处理mysql字段类型数据
     */
    public String dealColDtForMySQL(ResultSet rs, LoggerLog log, GetColumns gcl) throws SQLException {
        /**
         * define NULL for String constant
         */
        String dataStr = "NULL";
        for (String columnName : gcl.fieldmap.keySet()) {
            String dataType = gcl.fieldmap.get(columnName);
            // log.logging(dataType);
            if (dataType.equals("DATETIME")) {
                if (rs.getTimestamp(columnName, null) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr, rs.getTimestamp(columnName, null).toString());
                }

            } else {
                /**
                 * 其他数据类型：VARCHAR、TEXT等
                 */
                if (rs.getString(columnName) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr, rs.getString(columnName));

                    /* 遍历解析rule_josn */
                    if (this.tableName.equals("AM_RULE")) {
                        if (columnName.equals("rule_json")) {
                            // log.logging(dataStr);
                            try {
                                MyAmRule.parseRuleJson(rs.getString(columnName), log);
                            } catch (JsonSyntaxException e) {
                                log.logging(e.toString());
                            }
                        }

                        /* 提取风险标签 */
                        if (columnName.equals("label_id")) {
                            if (!rs.getString(columnName).isEmpty()) {
                                if (labelId.isEmpty()) {
                                    labelId = "'" + rs.getString(columnName) + "'";
                                } else {
                                    labelId = MyAmRule.appendValsSingle(labelId, rs.getString(columnName));
                                }
                            }
                        }
                    }
                    /*
                     * 嵌套指标提取:指标表
                     */
                    if (this.tableName.equals("AM_INDICATORS")) {
                        if (columnName.equals("indicators_param")) {
                            if (!rs.getString(columnName).isEmpty()) {
                                MyIndicator.parseIndicParam(rs.getString(columnName), log);
                            }
                        }
                    }
                    /*
                     * 提取指标模板ID
                     */
                    if (this.tableName.equals("AM_INDIC_TMPL")) {
                        if (columnName.equals("id")) {
                            if (indic_tmpl_ids.isEmpty() || indic_tmpl_ids == null) {
                                indic_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                indic_tmpl_ids = MyAmRule.appendValsSingle(indic_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }

                    /*
                     * 提取方法模板ID
                     */
                    if (this.tableName.equals("AM_METHOD")) {
                        if (columnName.equals("id")) {
                            if (mth_tmpl_ids.isEmpty() || mth_tmpl_ids == null) {
                                mth_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                mth_tmpl_ids = MyAmRule.appendValsSingle(mth_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }
                }
            }
        }
        return dataStr;
    }

    /*
     * 处理Oracle字段类型数据
     */
    public String dealColDtForOracle(ResultSet rs, LoggerLog log, GetColumns gcl)
            throws SQLException {
        /**
         * define NULL for String constant
         */
        String dataStr = "NULL";
        for (String columnName : gcl.fieldmap.keySet()) {
            String dataType = gcl.fieldmap.get(columnName);
            // 处理Oracle的数据类型：CLOB、BLOB、NCLOB、DATE等
            /**
             * 对DATE类型处理
             */
            if (dataType.equals("DATE")) {
                if (rs.getTimestamp(columnName, null) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr,
                            rs.getTimestamp(columnName, null).toString() + "00000");
                }
            }
            /**
             * 对NCLOB类型处理
             */
            if (dataType.equals("NCLOB")) {
                if (rs.getClob(columnName) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    try {
                        Clob clob = rs.getClob(columnName);
                        String clobString = this.ClobToString(clob);
                        dataStr = gcl.appendVals(dataStr, clobString);
                        /* 遍历解析rule_josn */
                        // log.logging(clobString);
                        if (this.tableName.equals("AM_RULE")) {
                            if (columnName.equals("RULE_JSON")) {
                                // log.logging(clobString);
                                MyAmRule.parseRuleJson(clobString, log);
                            }
                        }
                        /*
                         * 嵌套指标提取:指标表
                         */
                        if (this.tableName.equals("AM_INDICATORS")) {
                            if (columnName.equals("INDICATORS_PARAM")) {
                                MyIndicator.parseIndicParam(clobString, log);
                            }
                        }
                    } catch (JsonSyntaxException e) {
                        log.logging(e.toString());
                    } catch (SQLException e) {
                        log.logging(e.toString());
                    } catch (IOException e) {
                        log.logging(e.toString());
                    }

                }

            }
            /**
             * 对String类型处理
             */
            if (dataType.equals("NVARCHAR2")) {
                if (rs.getString(columnName) == null) {
                    dataStr = gcl.appendVals(dataStr, "NULL");
                } else {
                    dataStr = gcl.appendVals(dataStr, rs.getString(columnName));
                    /* 提取风险标签 */
                    if (columnName.equals("LABEL_ID")) {
                        if (!rs.getString(columnName).isEmpty()) {
                            if (labelId.isEmpty()) {
                                labelId = "'" + rs.getString(columnName) + "'";
                            } else {
                                labelId = MyAmRule.appendValsSingle(labelId, rs.getString(columnName));
                            }
                        }
                    }

                    /*
                     * 提取指标模板ID
                     */
                    if (this.tableName.equals("AM_INDIC_TMPL")) {
                        if (columnName.equals("ID")) {
                            if (indic_tmpl_ids.isEmpty() || indic_tmpl_ids == null) {
                                indic_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                indic_tmpl_ids = MyAmRule.appendValsSingle(indic_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }

                    /*
                     * 提取方法模板ID
                     */
                    if (this.tableName.equals("AM_METHOD")) {
                        if (columnName.equals("ID")) {
                            if (mth_tmpl_ids.isEmpty() || mth_tmpl_ids == null) {
                                mth_tmpl_ids = "'" + rs.getString(columnName) + "'";
                            } else {
                                mth_tmpl_ids = MyAmRule.appendValsSingle(mth_tmpl_ids, rs.getString(columnName));
                            }

                        }
                    }
                }
            }
        }
        return dataStr;
    }

    /**
     * CLOB转String
     */
    public String ClobToString(Clob clob) throws SQLException, IOException {
        String reString = "";
        Reader is = clob.getCharacterStream();// 得到流
        BufferedReader br = new BufferedReader(is);
        String s = br.readLine();
        StringBuffer sb = new StringBuffer();
        while (s != null) {
            // 执行循环将字符串全部取出付值给StringBuffer由StringBuffer转成STRING
            sb.append(s);
            s = br.readLine();
        }
        reString = sb.toString();
        return reString;
    }

}
