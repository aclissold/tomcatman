tomcatman
=========

**tomcatman**, the redundantly-named Tomcat multiple instance manager.

This program allows you to start, stop, and restart indivual Tomcat instances easily.

Usage
-----

First, ```go get github.com/aclissold/tomcatman```.

Change the line with ```os.Open("/path/to/tomcatman.conf")``` to a real path,
then ```go build tomcatman.go```.

Finally, configure tomcatman.conf. Add your instances in the following form:

    name=/path/to/tomcat1
    name=/path/to/tomcat2
    ...
    
where *name* is what you'd call each Tomcat instance and the path is its corresponding $CATALINA\_BASE.

Invoke the command with:

    tomcatman -start name
    tomcatman -stop name
    tomcatman -restart name

That's all there is to it.

Multiple Instance Setup
=======================

If you have not already set up multiple Tomcat instances but would like to, here's how.

The following is heavily based on 
[Running The Apache Tomcat 7.0 Servlet/JSP Container](http://tomcat.apache.org/tomcat-7.0-doc/RUNNING.txt);
visit it for a more detailed explanation.

General Idea
------------

In many circumstances, it is desirable to have a single copy of a Tomcat
binary distribution shared among multiple users on the same server. To make
this possible, you can set the CATALINA_BASE environment variable to the
directory that contains the files for your "personal" Tomcat instance.

If you already have a TOMCAT_HOME environment variable, it will become
CATALINA_HOME and files specific to each instance will be moved into a
new directory. Then, point CATALINA_BASE there when necessary.

Specifics
---------

The files and directories are split as following:

In CATALINA_BASE:

 * bin  - Only the following files:
          <br>&nbsp;&nbsp;&nbsp;&nbsp;* setenv.sh (*nix) or setenv.bat (Windows),
          <br>&nbsp;&nbsp;&nbsp;&nbsp;* tomcat-juli.jar
 * conf - Server configuration files (including server.xml)

 * lib  - Libraries and classes, as explained below

 * logs - Log and output files

 * webapps - Automatically loaded web applications

 * work - Temporary working directories for web applications

 * temp - Directory used by the JVM for temporary files (java.io.tmpdir)


In CATALINA_HOME:

 * bin  - Startup and shutdown scripts
          <br>&nbsp;&nbsp;&nbsp;&nbsp;The following files will be used only if they are absent in
          CATALINA_BASE/bin:
          <br>&nbsp;&nbsp;&nbsp;&nbsp;setenv.sh (*nix), setenv.bat (Windows), tomcat-juli.jar

 * lib  - Libraries and classes, as explained below

 * endorsed - Libraries that override standard "Endorsed Standards"
              libraries provided by JRE. See Classloading documentation
              in the User Guide for details.
              <br>&nbsp;&nbsp;&nbsp;&nbsp;By default this "endorsed" directory is absent.

In the default configuration the JAR libraries and classes both in
CATALINA_BASE/lib and in CATALINA_HOME/lib will be added to the common
classpath, but the ones in CATALINA_BASE will be added first and thus will
be searched first.

The idea is that you may leave the standard Tomcat libraries in
CATALINA_HOME/lib and add other ones such as database drivers into
CATALINA_BASE/lib.
