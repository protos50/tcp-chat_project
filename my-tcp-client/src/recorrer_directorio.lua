-- Función recursiva para recorrer un directorio
function recorrer_directorio(directorio, nivel, ignorar)
    local archivos = io.popen('ls -a "' .. directorio .. '"'):read('*all')
    for archivo in archivos:gmatch("[^\r\n]+") do
        if archivo ~= "." and archivo ~= ".." then
            local ruta_completa = directorio .. "/" .. archivo
            local atributos = io.popen('stat -c "%F" "' .. ruta_completa .. '"'):read('*all')
            local es_directorio = atributos:match("directory")

            if es_directorio then
                local relativa = string.rep("│   ", nivel - 1) .. "├── " .. archivo
                print(relativa)

                local ignorar_actual = false
                for _, patron in ipairs(ignorar) do
                    if archivo:match(patron) then
                        ignorar_actual = true
                        break
                    end
                end

                if not ignorar_actual then
                    recorrer_directorio(ruta_completa, nivel + 1, ignorar)
                end
            else
                local relativa = string.rep("│   ", nivel - 1) .. "└── " .. archivo
                print(relativa)
            end
        end
    end
end

-- Llama a la función con el directorio deseado y la lista de directorios a ignorar
local directorio_inicial = "/home/francozini/Documents/tcp-chat_project/my-tcp-client/src"
local directorios_a_ignorar = { ".git", "target" } -- Agrega los directorios que quieras ignorar
recorrer_directorio(directorio_inicial, 1, directorios_a_ignorar)
