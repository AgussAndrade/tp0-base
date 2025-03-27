# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

#### Solucion
##### protocolo cliente
El mensaje que representa una apuesta tiene el siguiente formato:

```
numero_cliente;nombre;apellido;dni;nacimiento;numero\n
```

- Campos separados por `;`
- Cada campo es un valor obligatorio y ordernado
- El mensaje termina con un **`\n`**

##### protocolo servidor
Segun si lo enviado fue correcto o no (puede ser un error debido a chequeo o un error porque llego alguna informacion mal).El servidor puede responder con:

```
OK\n
FAIL\n
```


---

##### Modificaciones realizadas en el cliente

- Se agrego una estructura `Bet` con los campos de la apuesta
- Se creo una funcion `getBets(id string) []Bet` que construye las apuestas desde variables de entorno
- Se cambio el bucle original para:
  - Enviar apuestas una a una
  - Armar manualmente el mensaje via metodo en utils
  - Enviar con bucle de `Write()` para evitar **short-write**
  - Leer byte a byte hasta `\n` para evitar **short-read**
  - Validar que la respuesta sea `"OK\n"`

##### Modificaciones realizadas en el servidor

- Se agregó una funcion `recv_until` para leer correctamente mensajes hasta el delimitador (evita **short-read**)
- Se utilizo la funcion `sendall` para enviar mensajes evitando **short-write**
- Se parseo el mensaje usando `split(';')` para obtener una lista de atributos ordenados para asi crear un `bet`
- Se devolvio la respuesta `"OK\n"` al cliente si fue exitosa
- Se devolvio la respuesta `"FAIL\n"` al cliente si fue erronea

##### Modificaciones realizadas en generar-compose.sh

Se modifico el script para agregar las variables de entorno de hasta 5 clientes leyendo de un archivo ubicado en `bet/{numero_client}.env` (en caso de querer agregar mas clientes hay que agregar su respectivo .env)

para correr el archivo ejecutar desde la raiz del proyecto

```bash
./generar-compose.sh docker-compose-dev.yaml 1
make docker-compose-up
```
