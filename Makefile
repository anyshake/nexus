.PHONY: lib clean

CFLAGS = -Wall -fPIC -I../
LDFLAGS = -shared

SRC = ../plugin.c
TARGET = libplugin.so

lib: $(TARGET)
$(TARGET): $(SRC)
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $^

clean:
	rm -f $(TARGET)
