import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.JsonSyntaxException;

public class MyIndicator {
    /*
     * 解析指标中的indic_params，
     * 提取主体：嵌套指标（来源：方法methods、过滤条件filterCondition）、模板、模板参数
     */
    public static void parseIndicParam(String indicators_param, LoggerLog log) throws JsonSyntaxException {
        JsonObject jsonObject = JsonParser.parseString(indicators_param).getAsJsonObject();
        JsonArray filterCondition = jsonObject.get("filterCondition").getAsJsonArray();
        JsonArray methods = jsonObject.get("methods").getAsJsonArray();
        if (filterCondition.size() != 0) {
            int i = 0;
            for (i = 0; i < filterCondition.size(); i++) {
                jsonObject = JsonParser.parseString(filterCondition.get(i).toString()).getAsJsonObject();
                // log.logging(jsonObject.toString());
                JsonArray parameters = jsonObject.get("parameters").getAsJsonArray();
                // log.logging(parameters.toString());
                if (parameters.size() != 0) {
                    int j = 0;
                    for (j = 0; j < parameters.size(); j++) {
                        jsonObject = JsonParser.parseString(parameters.get(j).toString()).getAsJsonObject();
                        jsonObject = jsonObject.get("value").getAsJsonObject();
                        String srcType = jsonObject.get("srcType").toString();
                        String srcId = jsonObject.get("id").toString();
                        String srcVal = jsonObject.get("value").toString();
                        srcType = srcType.substring(1, srcType.length() - 1);
                        srcId = srcId.substring(1, srcId.length() - 1);
                        srcVal = srcVal.substring(1, srcVal.length() - 1);
                        /**
                         * 开始提取主体：指标、名单、字段
                         */
                        if (srcType.equals("CALCULATION_VAR")) {
                            if (MyAmRule.calId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.calId = "'" + srcId + "'";
                            } else if (!MyAmRule.calId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.calId = MyAmRule.appendValsSingle(MyAmRule.calId, srcId);
                            }
                        }

                        if (srcType.equals("LIST_LIB")) {
                            if (MyAmRule.nameListId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.nameListId = "'" + srcId + "'";
                            } else if (!MyAmRule.nameListId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.nameListId = MyAmRule.appendValsSingle(MyAmRule.nameListId, srcId);
                            }
                        }
                        if (srcType.equals("BUSINESS_VAR")) {
                            if (MyAmRule.fieldId.isEmpty() && !srcVal.isEmpty()) {
                                MyAmRule.fieldId = "'" + srcVal + "'";
                            } else if (!MyAmRule.fieldId.isEmpty() && !srcVal.isEmpty()) {
                                MyAmRule.fieldId = MyAmRule.appendValsSingle(MyAmRule.fieldId, srcVal);
                            }
                        }
                        if (srcType.equals("ENUM_VAR")) {
                            if (MyAmRule.enumId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.enumId = "'" + srcId + "'";
                            } else if (!MyAmRule.enumId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.enumId = MyAmRule.appendValsSingle(MyAmRule.enumId, srcId);
                            }
                        }
                    }
                }

            }
        }
        /**
         * 解析方法methods中的主体
         */
        if (methods.size() != 0) {
            int i = 0;
            for (i = 0; i < methods.size(); i++) {
                jsonObject = JsonParser.parseString(methods.get(i).toString()).getAsJsonObject();
                // log.logging(jsonObject.toString());
                JsonArray parameters = jsonObject.get("parameters").getAsJsonArray();
                // log.logging("methods:" + parameters.toString());
                if (parameters.size() != 0) {
                    int j = 0;
                    for (j = 0; j < parameters.size(); j++) {
                        jsonObject = JsonParser.parseString(parameters.get(j).toString()).getAsJsonObject();
                        jsonObject = jsonObject.get("value").getAsJsonObject();
                        String srcType = jsonObject.get("srcType").toString();
                        srcType = srcType.substring(1, srcType.length() - 1);
                        /**
                         * 开始提取主体：指标、名单、字段
                         */
                        if (srcType.equals("CALCULATION_VAR")) {
                            String srcId = jsonObject.get("id").toString();
                            srcId = srcId.substring(1, srcId.length() - 1);
                            if (MyAmRule.calId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.calId = "'" + srcId + "'";
                            } else {
                                MyAmRule.calId = MyAmRule.appendValsSingle(MyAmRule.calId, srcId);
                            }
                        }
                        if (srcType.equals("LIST_LIB")) {
                            String srcId = jsonObject.get("id").toString();
                            srcId = srcId.substring(1, srcId.length() - 1);
                            // String srcVal = jsonObject.get("value").toString();
                            if (MyAmRule.nameListId.isEmpty() && !srcId.isEmpty()) {
                                MyAmRule.nameListId = "'" + srcId + "'";
                            } else {
                                MyAmRule.nameListId = MyAmRule.appendValsSingle(MyAmRule.nameListId, srcId);
                            }

                        }
                        if (srcType.equals("BUSINESS_VAR")) {
                            String srcVal = jsonObject.get("value").toString();
                            srcVal = srcVal.substring(1, srcVal.length() - 1);
                            if (MyAmRule.fieldId.isEmpty() && !srcVal.isEmpty()) {
                                MyAmRule.fieldId = "'" + srcVal + "'";
                            } else {
                                MyAmRule.fieldId = MyAmRule.appendValsSingle(MyAmRule.fieldId, srcVal);
                            }

                        }
                        if (srcType.equals("ENUM_VAR")) {
                            String srcVal = jsonObject.get("value").toString();
                            srcVal = srcVal.substring(1, srcVal.length() - 1);
                            if (MyAmRule.enumId.isEmpty() && !srcVal.isEmpty()) {
                                MyAmRule.enumId = "'" + srcVal + "'";
                            } else {
                                MyAmRule.enumId = MyAmRule.appendValsSingle(MyAmRule.enumId, srcVal);
                            }
                        }
                    }
                }

            }
        }
    }
}
