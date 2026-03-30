from getpass import getpass
import privy

data = getpass("Data: ")

print(privy.hide(data.encode("utf-8"), getpass("Password: ")))
