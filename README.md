# TP0: Docker + Comunicaciones + Concurrencia

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

#### Solucion

##### protocolo cliente  
Una vez que el cliente finaliza el envio de apuestas (con `\t`), permanece a la espera de un mensaje con los ganadores

- El cliente permanece a la espera de la respuesta por parte del servidor
- Enviar el final `\t` ya esta del ejercicio anterior como consecuencia de utilizar una unica conexion

La respuesta del servidor sera:

```
dni1;dni2;dni3;t
```

- Si no hubo ganadores, se responde simplemente con `\t`

##### protocolo servidor
El servidor ahora tambien maneja el proceso de sorteo. Una vez que todos los clientes hayan terminado de enviar sus apuestas, se ejecuta el sorteo

- El servidor mantiene el socket abierto con cada cliente luego del `\t`
- La respuesta es una lista de DNIs separados por `;`, terminada en `\t`
- Luego de enviar la respuesta, el servidor cierra el socket del cliente. Posible mejora: esperar un ok por parte de cada cliente

---

##### Modificaciones realizadas en el cliente

- Se extendio la logica para, luego de finalizar el envio de batchs, esperar respuesta del servidor (los ganadores)
- Se agrego parseo de la respuesta:
  - Se cuenta la cantidad de lineas (`split(';')`) hasta el `\t`
  - Se loguea `action: consulta_ganadores | result: success | cant_ganadores: X`
- Se mantiene la conexion abierta hasta recibir la respuesta y se cierra solo al final

##### Modificaciones realizadas en el servidor

- Se agrego una estructura interna para registrar que agencias ya finalizaron
- En el loop principal, ademas de esperar un sigterm, ahora se espera que lleguen todos los bets de los clientes propuestos por una variable de entorno seteada
- Una vez finalizado el loop principal, se calculan los ganadores para enviarlo a cada cliente en particular
- El servidor cierra el socket del cliente luego de enviar la lista de ganadores
- Se loguea `action: consulta_ganadores | result: success | cant_ganadores: X`
- En caso de fallo se loguea `action: consulta_ganadores | result: fail | agency: {agency_id}`

##### Modificaciones realizadas en generar-compose.sh

- Se agrego la variable de entorno `CLIENTNUMBER` en el servidor para saber cuantos clientes esperar antes de realizar el sorteo  

Para correr el archivo ejecutar desde la raiz del proyecto:

```bash
./generar-compose.sh docker-compose-dev.yaml 5
make docker-compose-up
```


