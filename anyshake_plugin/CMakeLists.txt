cmake_minimum_required(VERSION 3.10)
project(anyshake_plugin C)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -fPIC")
include_directories(../)

add_library(plugin STATIC ../plugin.c)

add_custom_command(TARGET plugin POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:plugin> ${CMAKE_CURRENT_BINARY_DIR}/libplugin.a
    COMMENT "Copying libplugin.a for Go CGO build"
)

set(GO_OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/anyshake_plugin)

add_custom_target(go-build
    COMMAND CGO_ENABLED=1 go build
        -ldflags=\"-s -w\"
        -v -trimpath
        -o ${GO_OUTPUT}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    DEPENDS plugin
    COMMENT "Building AnyShake plugin with CGO"
)

install(CODE "execute_process(
    COMMAND ${CMAKE_COMMAND} --build . --target go-build
    WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
)")
install(FILES ${GO_OUTPUT} DESTINATION ${SEEDLINK_PLUGIN_OUTPUT_DIR})

FILE(GLOB descs "${CMAKE_CURRENT_SOURCE_DIR}/descriptions/*.xml")
INSTALL(FILES ${descs} DESTINATION ${SC3_PACKAGE_APP_DESC_DIR})

add_custom_target(clean-all
    COMMAND ${CMAKE_COMMAND} -E remove -f ${GO_OUTPUT}
    COMMAND ${CMAKE_COMMAND} -E remove -f ${CMAKE_CURRENT_BINARY_DIR}/libplugin.a
    COMMENT "Cleaning up generated files"
)
