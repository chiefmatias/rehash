import socket
import sys

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 6379))

sock.sendall(b'$5\r\nHello\r\n')

#response = sock.recv(1024)
#print(response.decode('utf-8'))

sock.close()
sys.exit(0)
