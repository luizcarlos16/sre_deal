import os

def walk(path):
    f = []
    for root, dirs, files in os.walk(path, topdown=False):
       for name in files:
          f.append(os.path.join(root, name)) 
    return f
    
def check_if_string_in_file(file_name, string_to_search, string_to_replace):
    
    if string_to_search in open(file_name, 'r',encoding='latin-1').read():
        # print(file_name)
        newText = open(file_name, 'r',encoding='latin-1').read().replace(string_to_search, string_to_replace)
        # print(newText)
        with open(file_name, 'w',encoding='latin-1') as f:
            f.write(newText)

print('Qual diretorio quer analisar?')
caminho = input()
print('Qual string quer procurar?')
old = input()
print('Qual string quer adicionar?')
new = input()

for p in walk(caminho):
    check_if_string_in_file(p, old, new)