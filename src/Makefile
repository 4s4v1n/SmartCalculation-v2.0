BIN     = calculator
GOSRC   = main.go
CORE    = src/pkg/model/core
CXXSRC  = pkg/model/core/core_calculator.cpp pkg/model/core/core_creditor.cpp pkg/model/core/core_interface.cpp
CXXHDR  = pkg/model/core/core_calculator.hpp pkg/model/core/core_creditor.hpp pkg/model/core/core_interface.h
LIB     = libcalculator.so
SOFLAGS = -fPIC -shared
CXX     = g++
GO      = go
SUP     = --suppress=unusedFunction --suppress=missingIncludeSystem --suppress=unreadVariable
STYLE   = -style=Google

all: clean tidy cppcheck clang lib
	$(GO) build -o $(BIN) $(GOSRC)
	LD_LIBRARY_PATH=$(shell pwd)/pkg/model/core ./calculator

cppcheck:
	cppcheck --enable=all $(SUP) $(CXXSRC) $(CXXHDR)

clang:
	clang-format -n $(STYLE) $(CXXSRC) $(CXXHDR)

tidy:
	$(GO) mod tidy

lib:
	$(CXX) -o pkg/model/core/$(LIB) $(CXXSRC) $(SOFLAGS)

clean:
	rm -rf pkg/model/core/$(LIB) $(BIN)
