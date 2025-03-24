if [ -z "$1" ] || [ -z "$2" ]; then
    echo "use: $0 <file_name> <clients_number>"
    exit 1
fi

COMPOSE_FILE="$1"
CLIENT_NUMBER="$2"

cat <<EOF > "$COMPOSE_FILE"
name: tp0

networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24

services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    volumes:
      - ./server/config.ini:/config.ini
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net

EOF

for i in $(seq 1 "$CLIENT_NUMBER"); do
    ENV_FILE="./client/bet/${i}.env"

    if [ ! -f "$ENV_FILE" ]; then
        echo "$ENV_FILE not found. Skipping client$i"
        continue
    fi

    NOMBRE=$(grep '^NOMBRE=' "$ENV_FILE" | cut -d'=' -f2-)
    APELLIDO=$(grep '^APELLIDO=' "$ENV_FILE" | cut -d'=' -f2-)
    DOCUMENTO=$(grep '^DOCUMENTO=' "$ENV_FILE" | cut -d'=' -f2-)
    NACIMIENTO=$(grep '^NACIMIENTO=' "$ENV_FILE" | cut -d'=' -f2-)
    NUMERO=$(grep '^NUMERO=' "$ENV_FILE" | cut -d'=' -f2-)

    cat <<EOF >> "$COMPOSE_FILE"
  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    volumes:
      - ./client/config.yaml:/config.yaml
    environment:
      - CLI_ID=$i
      - NOMBRE=$NOMBRE
      - APELLIDO=$APELLIDO
      - DOCUMENTO=$DOCUMENTO
      - NACIMIENTO=$NACIMIENTO
      - NUMERO=$NUMERO
    networks:
      - testing_net
    depends_on:
      - server

EOF
done

echo "File $COMPOSE_FILE generated with $CLIENT_NUMBER clients."
