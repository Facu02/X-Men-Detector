
#### X-Men Detector
Este proyecto es una API en Go que permite detectar mutantes a partir de secuencias de ADN. Implementa dos endpoints: uno para obtener estadísticas y otro para validar si un conjunto de ADN corresponde a un mutante.

#### Requisitos
```
Go (versión 1.18 o superior)

Docker
```

#### Iniciar la Aplicación
Para ejecutar la aplicación en tu entorno local, debes seguir estos pasos:

1. Levantar Redis con Docker
Primero, asegúrate de tener Docker instalado en tu máquina. Luego, crea un archivo .env para almacenar las credenciales de tu base de datos Redis. Un archivo de ejemplo .env.example se proporciona en el proyecto.

```
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
```

En caso de estar ejecutando Redis de manera local, deja el campo __REDIS_PASSWORD__ vacío. Si estás usando una instancia de Redis remota, debes poner la contraseña correspondiente.

Luego, ejecuta el siguiente comando para levantar Redis en un contenedor Docker:

```
docker-compose up -d
```
Este comando levantará una instancia de Redis en un contenedor Docker. Asegúrate de tener el archivo docker-compose.yml configurado correctamente para usar Redis.

2. Crear un archivo .env
Crea un archivo .env en la raíz del proyecto con las siguientes variables de entorno (reemplaza los valores por los correspondientes en tu entorno local o remoto):

```
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
```

3. Iniciar la aplicación
Para iniciar la aplicación en tu entorno local, corre el siguiente comando:

```
go run main.go
```
La API se ejecutará en localhost:8080.

### Endpoints

#### GET /stats
Obtén estadísticas sobre el número de mutantes y humanos.

localhost:8080/stats

Respuesta de ejemplo:

```json
{
    "mutant_count": 1,
    "human_count": 0,
    "ratio": 1
}
```
Esta respuesta muestra:

mutant_count: La cantidad de mutantes detectados.
human_count: La cantidad de humanos detectados.
ratio: El ratio de mutantes respecto a humanos.
### POST /mutant
Verifica si el ADN recibido corresponde a un mutante. El cuerpo del mensaje debe contener una matriz de ADN en formato JSON. La matriz debe ser de al menos 4x4.

localhost:8080/mutant

Cuerpo de ejemplo:

```json
{
    "dna": [
        "ATGCGA",
        "CAGTGC",
        "TTATGT",
        "AGAAGG",
        "CCCCTA",
        "TCACTG"
    ]
}
```

Si el ADN corresponde a un mutante, la respuesta será 
``{"is_mutant": true}.``
En caso contrario, será 
403 forbidden.

