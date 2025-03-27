# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

#### Solucion

##### protocolo cliente
El mensaje ahora representa **varias apuestas agrupadas** en un solo envio (batch), con el siguiente formato:

```
numero_cliente;nombre;apellido;dni;nacimiento;numero\n
numero_cliente;nombre;apellido;dni;nacimiento;numero\n
\t
```

- Cada apuesta es un mensaje ordenado separado por `\n`
- El batch completo termina con **`\t`**
- El cliente puede enviar **multiples batchs** por conexion
- El cliente reintentara hasta 3 veces si un batch falla desde el server. si se alcanza el maximo se cortara la ejecucion
- El cliente envia un mensaje vacio con `\t` para indicar que finalizo el envio

##### protocolo servidor
El servidor ahora valida **todo el batch recibido**. Segun el resultado, responde con:

```
OK\t
FAIL\t
```

- Si **todas** las apuestas del batch son validas, responde `OK\t`
- Si **alguna** apuesta es invalida, responde `FAIL\t` y descarta el batch completo
- El servidor permite reintentos por parte del cliente
- La conexion permanece abierta para múltiples batchs

---

##### Modificaciones realizadas en el cliente

- Se agrego una logica para leer apuestas desde un archivo CSV (`agency-{id}.csv`)
- Se agruparon apuestas en batchs respetando:
  - Un limite maximo de apuestas por batch (`batch_maxAmount`) seteado en el config.ini
  - Un limite maximo en bytes (`batch_maxBytes`) seteado en el docker-compose
- Se implemento una funcion `retryBatchUntilSuccess` que reintenta el envio del batch si recibe `FAIL`
- Se modifico el envio para:
  - Enviar cada batch como un string con `\n` entre apuestas, terminado en `\t`
  - Leer la respuesta byte a byte hasta `\t` (evita **short-read**)
- Se envia un `\t` solo al final para indicar fin de transmisión

##### Modificaciones realizadas en el servidor

- Se adapto el protocolo de recepcion para leer hasta `\t` usando `recv_until`
- Se parseo el mensaje como una lista de lineas (`split('\n')`) y luego cada linea como una apuesta (`split(';')`)
- Se valido el batch completo:
  - Si todas son validas se almacena con `store_bets` y se responde `OK\t`
  - Si alguna es invalida se responde `FAIL\t` y se permite el reenvio
- Se mantuvo la conexion abierta para recibir multiples batchs por socket

##### Modificaciones realizadas en generar-compose.sh

- se agrego el volumen para el archivo `.data/agency-{n}.csv` para que cada cliente tenga acceso a su archivo de apuestas
- Se borraron las logicas para obtener las apuestas con respecto al ejercicio anterior

para correr el archivo ejecutar desde la raiz del proyecto:

```bash
./generar-compose.sh docker-compose-dev.yaml 1
make docker-compose-up
```
