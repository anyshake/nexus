cmake_minimum_required(VERSION 3.10)
project(anyshake_plugin NONE)

set(GO_OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/anyshake_plugin)

add_custom_target(go-build
    COMMAND ${CMAKE_COMMAND} -E echo "Building Go plugin"
    COMMAND CGO_ENABLED=0 go build
        -ldflags=-s\ -w
        -v -trimpath
        -o ${GO_OUTPUT}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    COMMENT "Building AnyShake plugin (zero dependency)"
)

add_custom_target(go-build-target ALL
    DEPENDS go-build
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
    COMMENT "Cleaning up Go binary"
)
