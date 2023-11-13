javac -encoding utf-8 -d classes\ HelloWorld.java
jar cvmf manifest.mf HelloWorld.jar -C classes\ .
java -jar HelloWorld.jar



javadoc -encoding utf-8 -d docs -author zh -version JavaDoc.java