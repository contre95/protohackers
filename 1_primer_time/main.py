import sys

def is_bip39_word(word):
    with open('bip39_wordlist.txt', 'r') as f:
        bip39_words = {line.strip() for line in f}
    return word in bip39_words

input_string = sys.stdin.readline().strip()

a = {w for w in input_string.split(" ") if is_bip39_word(w)}

for w in a:
    print(w, end=" ")
print()
