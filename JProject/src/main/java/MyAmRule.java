import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.JsonSyntaxException;

public class MyAmRule {
    public static String calId = "";
    public static String nameListId = "";
    public static String fieldId = "";
    public static String enumId = "";

    /* 解析规则中的rule_json，提取主体：字段、指标、名单 */
    public static void parseRuleJson(String rule_json, LoggerLog log) throws JsonSyntaxException {
        JsonObject jsonObject = JsonParser.parseString(rule_json).getAsJsonObject();
        JsonArray conditionList = jsonObject.get("conditionList").getAsJsonArray();
        // log.logging(conditionList.get(0).toString());
        jsonObject = JsonParser.parseString(conditionList.get(0).toString()).getAsJsonObject();
        JsonArray conditionItems = jsonObject.get("conditionItems").getAsJsonArray();
        /**
         * 遍历条件主体，获取fieldSource，fieldValue；
         * 根据fieldSource分类，统计收集fieldValue（应去重），到对应的实体表查询
         */
        JsonObject targetJson;
        String targetSource, targetValue;
        int i = 0;
        for (i = 0; i < conditionItems.size(); i++) {
            targetJson = JsonParser.parseString(conditionItems.get(i).toString()).getAsJsonObject();
            targetSource = targetJson.get("fieldSource").toString();
            targetValue = targetJson.get("fieldValue").toString();
            targetSource = targetSource.substring(1, targetSource.length() - 1);
            targetValue = targetValue.substring(1, targetValue.length() - 1);
            /**
             * 主体区分：指标、业务字段、名单
             */
            if (targetSource.equals("CALCULATION_VAR")) {
                if (calId.isEmpty() && !targetValue.isEmpty()) {
                    calId = "'" + targetValue + "'";
                } else if (!calId.isEmpty() && !targetValue.isEmpty()) {
                    calId = MyAmRule.appendValsSingle(calId, targetValue);
                }
                // log.logging("Source:" + targetSource + "; Value:" + targetValue);
            }
            if (targetSource.equals("BUSINESS_VAR")) {
                if (fieldId.isEmpty() && !targetValue.isEmpty()) {
                    fieldId = "'" + targetValue + "'";
                } else if (!fieldId.isEmpty() && !targetValue.isEmpty()) {
                    fieldId = MyAmRule.appendValsSingle(fieldId, targetValue);
                }
                // log.logging("Source:" + targetSource + "; Value:" + targetValue);
            }
            /*
             * 枚举字段
             */
            if (targetSource.equals("ENUM_VAR")) {
                if (enumId.isEmpty() && !targetValue.isEmpty()) {
                    enumId = "'" + targetValue + "'";
                } else if (!enumId.isEmpty() && !targetValue.isEmpty()) {
                    enumId = MyAmRule.appendValsSingle(enumId, targetValue);
                }
                // log.logging("Source:" + targetSource + "; Value:" + targetValue);
            }
            if (targetSource.equals("LIST_LIB")) {
                if (nameListId.isEmpty() && !targetValue.isEmpty()) {
                    nameListId = "'" + targetValue + "'";
                } else if (!nameListId.isEmpty() && !targetValue.isEmpty()) {
                    nameListId = MyAmRule.appendValsSingle(nameListId, targetValue);
                }
                // log.logging("Source:" + targetSource + "; Value:" + targetValue);
            }
        }
    }

    // 将字段拼接成字符串，处理字段为空输出逻辑+字段拼接逻辑:去重
    public static String appendValsSingle(String rec, String val) {
        if (!val.equals("NULL")) {
            val = "'" + val + "'";
        }
        // 去重统计
        if (!rec.contains(val) && !val.isEmpty()) {
            if (!rec.equals("NULL")) {
                rec = rec + "," + val;
            } else {
                rec = val;
            }
        }
        return rec;
    }
}
