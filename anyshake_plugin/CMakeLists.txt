cmake_minimum_required(VERSION 3.10)
project(anyshake_plugin C)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -fPIC")
include_directories(../../libs/plugin)

add_library(plugin STATIC ../../libs/plugin/plugin.c)

add_custom_command(TARGET plugin POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:plugin> ${CMAKE_CURRENT_SOURCE_DIR}/libplugin.a
    COMMENT "Copying libplugin.a for Go CGO build"
)

file(GLOB GO_SOURCES "${CMAKE_CURRENT_SOURCE_DIR}/*.go")

set(GO_OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/anyshake_plugin)

add_custom_command(
    OUTPUT ${GO_OUTPUT}
    COMMAND CGO_ENABLED=1 go build
        -ldflags=\"-s -w\"
        -v -trimpath
        -o ${GO_OUTPUT}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    DEPENDS plugin ${GO_SOURCES}
    COMMENT "Building AnyShake plugin with CGO"
)

add_custom_target(go-build-target ALL
    DEPENDS ${GO_OUTPUT}
)

install(FILES ${GO_OUTPUT}
    DESTINATION ${SEEDLINK_PLUGIN_OUTPUT_DIR}
    PERMISSIONS OWNER_READ OWNER_WRITE OWNER_EXECUTE
                GROUP_READ GROUP_EXECUTE
                WORLD_READ WORLD_EXECUTE
)

file(GLOB descs "${CMAKE_CURRENT_SOURCE_DIR}/descriptions/*.xml")
install(FILES ${descs} DESTINATION ${SC3_PACKAGE_APP_DESC_DIR})

add_custom_target(clean-all
    COMMAND ${CMAKE_COMMAND} -E remove -f ${GO_OUTPUT}
    COMMAND ${CMAKE_COMMAND} -E remove -f ${CMAKE_CURRENT_SOURCE_DIR}/libplugin.a
    COMMENT "Cleaning up generated files"
)
