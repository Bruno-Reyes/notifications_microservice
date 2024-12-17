# Instrucciones generales

## Inicializar un módulo

```bash
go mod init <project-name>
```

Esto crea un archivo go.mod que administrará las dependencias de tu proyecto.

### Instalación de librerias
```bash
go get <librarie-name>
```
* Descarga la librería.
* Actualiza el archivo go.mod y go.sum con la información del módulo.

## Para descargar las dependencias 

Similar a "bun run dev" ó "pip install -r requirements.txt"
```bash
go mod tidy
```
Esto asegura que todas las dependencias necesarias estén instaladas y elimina las que no se usan.

# Ejecucion del proyecto
```bash
go run main.go
```