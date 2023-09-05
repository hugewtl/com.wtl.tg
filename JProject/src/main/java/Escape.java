public class Escape {
	/*
	 * 针对特殊字符进行转义：
	 */
	public String escapeForSC(String str) {
		/*
		 * '%' 比较特殊，放在第一位，如果原字符串包含%,则进行转义
		 */
		String[] specialChar = new String[] { "%", "!", "@", "#", "$", "^", "&", "*", "(", ")", "_", "+", "=" };
		for (String substr : specialChar) {
			if (str.contains(substr)) {
				str = str.replace(substr, "\\" + substr);
				System.out.println(str);
			}
		}
		return str;
	}
}
/*
 * 转义密码中的特殊字符
 */
// Escape escape = new Escape();
// dbi.setDbPasswd(escape.escapeForSC(rp.properties.getProperty("dbpasswd")));