# 
#	vgof 框架
#	

# 主程序
MAIN_PROG = ican
# 模块
MODULES = configMod icanMod sendmailMod zaplogMod
# 编译选项 -s -w
LDFLAGS = 

.PHONY: all $(MAIN_PROG) $(MODULES)

all: $(MAIN_PROG) $(MODULES)

$(MAIN_PROG): 
	go build -ldflags "$(LDFLAGS)" -o $@

$(MODULES): 
	cd ./modules/$@ && go build -buildmode=plugin -ldflags "$(LDFLAGS)" && \
	mv -f ./$@.so ../$@.so
