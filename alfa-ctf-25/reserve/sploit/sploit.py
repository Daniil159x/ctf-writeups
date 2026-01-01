from scapy.all import *
import sys

def xor_payload_with_key(payload, key):
    key_len = len(key)
    return bytes([payload[i] ^ key[i % key_len] for i in range(len(payload))])

def process_pcap(input_file, output_file):
    packets = rdpcap(input_file)
    decrypted = []

    MASK = (1 << 32) - 1
    for pkt in packets:
        if TCP in pkt and pkt[TCP].sport == 1337:
            tcp_payload = bytes(pkt[TCP].payload)

            if len(tcp_payload) > 4:
                key = int.from_bytes(pkt[Padding], 'little')

                origin_payload = []
                for i in range(0, len(tcp_payload), 4):
                    crypted = int.from_bytes(tcp_payload[i:i+4], 'little')
                    plain = crypted ^ key

                    key = (0x88767bc1 * key) & MASK
                    key ^= crypted
                    key = (0x20763841 * key) & MASK

                    origin_payload += plain.to_bytes(4, 'little')

                origin_payload = origin_payload[:len(tcp_payload) - 4]

                # Modify payload
                pkt[TCP].remove_payload()
                pkt[TCP].add_payload(bytes(origin_payload))

            # delete key
            if Padding in pkt:
                del pkt[Padding]
        print(pkt)
        decrypted.append(pkt)

    # Save new pcap
    wrpcap(output_file, decrypted)

# Example usage
input_pcap = sys.argv[1]
output_pcap = sys.argv[2]
process_pcap(input_pcap, output_pcap)