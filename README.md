pyx是一个使用go语言编写的工具，用于准备好启动一个docker容器的命令并执行它。

pyx的代码逻辑就是创建一个本地临时目录并映射成容器中的/deps-cache目录，将当前项目目录映射成容器中的/build目录，将项目目录下的deps.txt文件中的目录路径映射到容器中的/private-cache目录下（比如/Users/mac/common/utils会被映射为/private-cache/utils，depts.txt文件中可以放置多个目录路径，一行一个路径），并将容器中的/private-cache设置到PYTHONPATH环境变量中，然后使用docker run --rm启动一个容器，默认启动的镜像名称是x86_64_pyx（如果必要，你可以在代码中找到它并修改它。）

按照如下步骤安装pyx工具：

1、安装go解释器，下载地址

https://studygolang.com/dl

2、将pyx.go下载到你的目录下，然后切换到你的目录下，执行如下命令打包并安装它

```
go install pyx.go
```

执行后，pyx会被安装到你的~/go/bin目录中，比如我的就是：/Users/mac/go/bin

需要将~/go/bin添加到PATH环境中，以便在命令行执行pyx时，可以找到这个pyx程序。



我还编写了一个Dockerfile文件，用于打包出一个镜像。

Dockerfile的打包逻辑是使用基础镜像centos:7，首先安装了V3.6.8版本的python，并将python可执行文件的路径添加到PATH环境变量中，并使用pip安装了pip-download，并将/build作为工作目录。镜像被作为容器启动时，会执行如下的命令：

（1）在当前项目目录下创建./build/stable目录，并删除./build/stable目录下除了libs目录以外的所有文件。

（2）使用pip-download将当前项目目录下的requirements.txt中的依赖下载到./build/stable/libs目录下。

（3）将当前项目目录下除了./build目录外的所有文件拷贝到./build/stable/目录下。

（4）将容器中/private-cache目录中的所有文件拷贝到./build/stable/目录下。

（5）将./build/stable目录压缩成./build/stable.tar。

（6）在容器中的/deps-cache/目录下中创建venv虚拟环境，名称为pyx。

（7）激活activate虚拟环境，并在虚拟环境中安装cx_Freeze。

（8）在虚拟环境中安装requirements.txt中的所有依赖。

（9）使用python setup.py build命令打包可执行文件。

（8）将打包后的./build/exe.linux-x86_64-3.6目录，压缩为./build/exe.linux-x86_64-3.6.tar文件。



按照如下步骤打包镜像：

1、首先需要在你的本地安装上docker，具体可以参考网上的教程。

2、进入Dockerfile所在的目录下，构建镜像

```
docker pull centos:7
docker build -t x86_64_pyx .
```



将pyx工具安装好，并且将docker镜像构建完成后，你就可以愉快的使用它进行项目打包了。

使用方法为：

1、首先遵循cx_Freeze打包的要求，在你的项目目录下创建setup.py文件。

2、在你的项目目录下创建一个deps.txt文件，将你的本地公共代码的路径添加到里面。

PS：如果没有创建deps.txt文件也没关系，打包能够正常执行。

3、然后打开终端命令行，进入到你的项目目录下执行打包命令

```
pyx
```



Yes，打包就是如此简单，只需要执行pyx即可，无任何参数。