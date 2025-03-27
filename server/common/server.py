import socket
import logging
import signal
import sys

from .utils import *

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True

        signal.signal(signal.SIGTERM, self._graceful_shutdown)

    def _graceful_shutdown(self, signum, term):
        logging.info(f'action: shutdown | result: in_progress | signal: {signum}')
        self._running = False
        try:
            self._server_socket.close()
            logging.info('action: close_socket | result: success')
        except Exception as e:
            logging.error(f'action: close_socket | result: fail | error: {e}')
        sys.exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    self.__handle_client_connection(client_sock)
            except OSError:
                break

    def __handle_client_connection(self, client_sock):
        try:
            while True:
                raw_data = recv_until(client_sock, delimiter=b'\t')
                if raw_data == b'':
                    logging.warning("action: receive_batch | result: fail | reason: timeout or empty")
                    break

                raw_data = raw_data.strip().decode('utf-8').strip('\t')
                if raw_data == '':
                    logging.info("action: receive_batch | result: success | reason: empty batch")
                    break

                lines = [line for line in raw_data.split('\n') if line.strip()]

                bets = []
                for line in lines:
                    msg = line.strip().split(';')
                    if not check_bet_msg(msg):
                        client_sock.sendall(b'FAIL\t')
                        logging.warning(f'action: apuesta_recibida | result: fail | cantidad: {len(lines)} | msg: {msg}')
                        break
                    bet = construct_bet_by_msg(msg)
                    bets.append(bet)
                
                store_bets(bets)
                client_sock.sendall(b'OK\t')
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                continue

        except Exception as e:
            logging.error(f'action: receive_batch | result: fail | error: {e}')
            client_sock.sendall(b'FAIL\t')
        finally:
            client_sock.close()
            logging.info('action: close_client_socket | result: success')

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c
        except OSError:
            return None
