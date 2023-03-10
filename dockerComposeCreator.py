import sys

def create(x):
	file=open("Docker-compose.yaml","w")
	file.write("version: \"3.9\"\n")
	file.write("services:\n")
	file.write("  registration_server:\n")
	file.write("    hostname: register\n")
	file.write("    container_name: register\n")
	file.write("    build:\n")
	file.write("      context: .\n")
	file.write("      dockerfile: ./DockerFile/server_register\n")
	file.write("    ports:\n")
	file.write("      - 8000:8000\n")
	file.write("    networks:\n")
	file.write("      - mynetwork\n")
	for i in range(1,x+1):
		s=str(i)
		file.write("  nodo"+s+":\n")
		file.write("    container_name: nodo"+s+"\n")
		file.write("    hostname: nodo"+s+"\n")
		file.write("    build:\n")
		file.write("      context: .\n")
		file.write("      dockerfile: ./DockerFile/nodo\n")
		file.write("    ports:\n")
		file.write("      - 800"+s+":8005\n")
		file.write("    environment:\n")
		file.write("      - PORT_EXSPOST=800"+s+"\n")
		file.write("    depends_on:\n")
		file.write("      - registration_server\n")
		file.write("    networks:\n")
		file.write("      - mynetwork\n")
	file.write("networks:\n")
	file.write("  mynetwork:\n")
	file.write("   name: network1234\n")
z=int(sys.argv[1])
create(z)