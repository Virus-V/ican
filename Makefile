# 
#	vgof 框架
#	

# 主程序
MAIN_PROG = ican
# 模块
MODULES = configMod icanMod sendmailMod zaplogMod

.PHONY: all $(MAIN_PROG) $(MODULES)

all: $(MAIN_PROG) $(MODULES)

$(MAIN_PROG): 
	go build -o $@

$(MODULES): 
	cd ./modules/$@ && go build -buildmode=plugin && \
	mv -f ./$@.so ../$@.so
