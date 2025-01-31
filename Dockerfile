FROM centos:7
RUN rm -rf /etc/yum.repos.d/* && curl -o /etc/yum.repos.d/CentOS7-Aliyun.repo http://mirrors.aliyun.com/repo/Centos-7.repo
RUN yum -y install wget xz unzip gcc make zlib zlib-devel openssl openssl-devel libffi-devel
RUN wget https://mirrors.huaweicloud.com/python/3.6.8/Python-3.6.8.tar.xz
RUN tar -xvf Python-3.6.8.tar.xz \
    && cd Python-3.6.8 \
    && ./configure --prefix=/usr/local/python --enable-shared \
    && make && make install
RUN echo 'export PATH=/usr/local/python/bin:$PATH' >> /etc/profile \
    && echo 'export LD_LIBRARY_PATH=/usr/local/python/lib:$LD_LIBRARY_PATH' >> /etc/profile
RUN mkdir -p /build /deps-cache /private-cache ~/.config/pip/
RUN echo '[global]' >> ~/.config/pip/pip.conf \
    echo 'index-url = https://pypi.tuna.tsinghua.edu.cn/simple' >> ~/.config/pip/pip.conf \
    echo '[install]' >> ~/.config/pip/pip.conf \
    echo 'trusted-host = pypi.tuna.tsinghua.edu.cn' >> ~/.config/pip/pip.conf
RUN source /etc/profile && pip3 install pip -U && pip3 install pip-download
WORKDIR /build
CMD source /etc/profile \
    && mkdir -p ./build/stable && rm -rf ./build/stable.tar \
    && find ./build/stable/ -mindepth 1 -maxdepth 1 -not -name 'libs' -exec rm -rf {} \; \
    && pip-download -i https://mirrors.aliyun.com/pypi/simple/ -r requirements.txt -d ./build/stable/libs \
    && echo "copying source code in current directory to ./build/stable/ ..." \
    && find . -mindepth 1 -maxdepth 1 -not -name 'build' -exec cp -r {} ./build/stable/ \; \
    && echo "copying source code in deps.txt to ./build/stable/ ..." \
    && ([ "$(ls -A /private-cache)" ] && cp -r /private-cache/* ./build/stable/) \
    && echo "starting compress ./build/stable.tar ..." \
    && cd ./build && tar -cf stable.tar stable && cd .. \
    && python3 -m venv /deps-cache/pyx \
    && source /deps-cache/pyx/bin/activate \
    && echo "starting package using cx_Freeze ..." \
    && pip install pip -U && pip install cx_Freeze \
    && pip install -r requirements.txt --no-index --find-links=./build/stable/libs \
    && python setup.py build \
    && cd ./build && && rm -rf exe.linux-x86_64-3.6.tar \
    && tar -cf exe.linux-x86_64-3.6.tar exe.linux-x86_64-3.6 \
    && echo "package successfully!"

# 类似于xgo通过docker使用cx_freeze自动给python代码打包linux版本的可执行文件