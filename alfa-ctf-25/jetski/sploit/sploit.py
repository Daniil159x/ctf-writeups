#!/usr/bin/env python3

import logging
import socket
from telnetlib import DO

def parse_pair(data: str):
    x, y = data[1:-1].split(',')
    return int(x), int(y)

def parse_binding(binding):
    keys = binding[0::2]
    dirs = binding[1::2]
    return {d:k for k, d in zip(keys, dirs)}

LEFT = 'l'
RIGHT = 'r'
UP = 'u'
DOWN = 'd'
MAX_PER_PING = 7

class Game:

    def __init__(self, host, port) -> None:
        self._sock = socket.socket()
        self._sock.connect((host, port))

        self.binding = parse_binding("WuAlSdDr")

    def start(self, token:str):
        self._sock.send(token.encode())

        init_state = self._recv_message()

        # {8,4}|{65,32}|[{42,14}+{56,29}+{62,14}+{64,22}]
        player, sizes, stars_arr = init_state.split('|')
        player = parse_pair(player)
        sizes = parse_pair(sizes)
        # stars = [parse_pair(star) for star in stars_arr.split('+')]

        return player, sizes

    def turns(self, dirs):
        pressed = ""
        for d in dirs:
            if len(pressed) == MAX_PER_PING:
                player, stars, _, msg = self.ping(pressed)
                pressed = ""

                if msg != "" or stars == 0:
                    raise Exception(f"Game over: {msg=} {stars=}")
            pressed += self.binding[d]

        player, stars, _, msg = self.ping(pressed)
        if msg != "" or stars == 0:
            raise Exception(f"Game over: {msg=} {stars=}")

        return player, stars

    def ping(self, pressed):
        pressed = '+'.join(pressed)
        # [d++d++++d]
        self._send_message(f"[{pressed}]")
        
        pong = self._recv_message()

        # {19,4}|293|WuAdSlDr|somemsg
        player, stars_left, binding, msg = pong.split('|')
        player = parse_pair(player)
        stars_left = int(stars_left)
        self.binding = parse_binding(binding)

        return player, stars_left, self.binding, msg

    def _recv_message(self):
        buf = b''
        while True:
            buf += self._sock.recv(4096)
            if buf.endswith(b"kek\n"):
                break
        if not buf.startswith(b"lol"):
            raise Exception(f"unexpected msg: {buf}")

        msg = buf[3:-4].decode()
        logging.debug(f"read msg: {msg}")
        return msg
    
    def _send_message(self, msg):
        logging.debug(f"send msg: {msg}")
        self._sock.send(f"lol{msg}kek\n".encode())

def turn(binding, dir, n):
    return [binding[dir] for _ in range(n)]

def main(host, port, token):
    server = Game(host, port)
    player, sizes = server.start(token)

    # просто идёт в верхний левый угол, а затем построчно пылесосим
    logging.info(f"connected {player=} {sizes=}")

    player, stars = server.turns([LEFT] * player[0])
    logging.info(f"turned left {player=} {stars=}")

    player, stars = server.turns([UP] * player[1])
    logging.info(f"turned up {player=} {stars=}")

    for y in range(0, sizes[1], 2):
        player, stars = server.turns([RIGHT] * sizes[0])
        logging.info(f"{y=} line done {player=} {stars=}")

        if y + 1 == sizes[1]:
            break

        player, stars = server.turns([DOWN])
        player, stars = server.turns([LEFT] * sizes[0])
        logging.info(f"{y+1=} line done {player=} {stars=}")

        if y + 2 == sizes[1]:
            break
        player, stars = server.turns([DOWN])

    logging.info(f"vacuuming done {player=} {stars=}")


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="jetski-455wf7wv.alfactf.ru")
    parser.add_argument("--port", default="30036")
    parser.add_argument("--token", default="wfewfewfwefweewf")

    logging.basicConfig(level=logging.INFO)
    try:
        args = parser.parse_args()
        main(args.host, int(args.port), args.token)
    except Exception as e:
        if "alfa{" in str(e):
            logging.warning(f"FLAG: {str(e)}")
        else:
            raise e
