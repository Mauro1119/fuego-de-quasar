package main

import (
	"fmt"
	"math"
	"strings"
)

//GetLocation obtiene la localizacion del punto emisor del mensaje.
func GetLocation(distances ...float32) (x, y float32) {
	//Verficar si se tocan las circunferencias sino error.
	var d1 float32
	var d2 float32
	var d3 float32
	sigueOk := true

	d1 = 670.8203  //Raiz de 450000
	d2 = 447.2135  //Raiz de 200000
	d3 = 1044.0306 //Raiz de 1090000

	if !((d1 <= distances[0]+distances[1]) && (d1 >= float32(math.Abs(float64(distances[0]-distances[1]))))) {
		sigueOk = false
	}

	if !((d2 <= distances[1]+distances[2]) && (d2 >= float32(math.Abs(float64(distances[1]-distances[2]))))) {
		sigueOk = false
	}

	if !((d3 <= distances[0]+distances[2]) && (d3 >= float32(math.Abs(float64(distances[0]-distances[2]))))) {
		sigueOk = false
	}

	if !sigueOk {
		return 100, -100
	}

	var M [2][3]float32

	D := (-2) * (-500)
	E := (-2) * (-200)
	F := (-500)*(-500) + (-200)*(-200) - distances[0]*distances[0]

	M[0][0] = float32(D)
	M[0][1] = float32(E)
	M[0][2] = (-1) * float32(F)

	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
	fmt.Println("")

	D = (-2) * (100)
	E = (-2) * (-100)
	F = (100)*(100) + (-100)*(-100) - distances[1]*distances[1]

	M[0][0] -= float32(D)
	M[0][1] -= float32(E)
	M[0][2] -= (-1) * float32(F)

	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
	fmt.Println("")

	D = (-2) * (-500)
	E = (-2) * (-200)
	F = (-500)*(-500) + (-200)*(-200) - distances[0]*distances[0]

	M[1][0] = float32(D)
	M[1][1] = float32(E)
	M[1][2] = (-1) * float32(F)

	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
	fmt.Println("")

	D = (-2) * (500)
	E = (-2) * (100)
	F = (500)*(500) + (100)*(100) - distances[2]*distances[2]

	M[1][0] -= float32(D)
	M[1][1] -= float32(E)
	M[1][2] -= (-1) * float32(F)

	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
	fmt.Println("")

	if M[0][0] != 0 {
		multi := M[0][0]
		for i := 0; i < 3; i++ {
			M[0][i] = M[0][i] / multi
		}
		fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
		fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
		fmt.Println("")
	} else {
		if M[1][0] != 0 {
			var cambio float32
			for i := 0; i < 3; i++ {
				cambio = M[0][i]
				M[0][i] = M[1][i]
				M[1][i] = cambio
			}

			fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
			fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
			fmt.Println("")

			multi := M[0][0]
			for i := 0; i < 3; i++ {
				M[0][i] = M[0][i] / multi
			}

			fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
			fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
			fmt.Println("")
		} else {
			return 100, -100
		}
	}

	multi := M[1][0]
	for i := 0; i < 3; i++ {
		M[1][i] = M[1][i] - M[0][i]*multi
	}

	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
	fmt.Println("")

	if M[1][1] != 0 {
		multi := M[1][1]
		for i := 1; i < 3; i++ {
			M[1][i] = M[1][i] / multi
		}

		fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
		fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")
		fmt.Println("")
	} else {
		return 100, -100
	}

	multi = M[0][1]
	for i := 1; i < 3; i++ {
		M[0][i] = M[0][i] - M[1][i]*multi
	}
	fmt.Println("[", M[0][0], " ", M[0][1], " ", M[0][2], "]")
	fmt.Println("[", M[1][0], " ", M[1][1], " ", M[1][2], "]")

	return M[0][2], M[1][2]
}

//GetMessage descifra el mensaje.
func GetMessage(messages ...[]string) (msg string) {
	var largo int

	for y := 0; y < len(messages); y++ {
		if y == 0 {
			largo = len(messages[y])
		} else {
			if largo > len(messages[y]) {
				largo = len(messages[y])
			}
		}
	}

	mensaje := []string{}
	var noEncontro bool
	for l := 0; l < largo; l++ {
		noEncontro = true
		for _, m := range messages {

			if largo < len(m) {
				if m[l+(len(m)-largo)] != "" {
					noEncontro = false
					mensaje = append(mensaje, m[l+len(m)-largo])
					break
				}
			} else {
				if m[l] != "" {
					noEncontro = false
					mensaje = append(mensaje, m[l])
					break
				}
			}

		}
		if noEncontro {
			return ""
		}
	}
	return strings.Join(mensaje, " ")

}

//Find busca el satelite.
func Find(sats []satellite, nombre string) (dist float32) {
	for _, d := range sat.Satellite {
		if strings.ToUpper(d.Name) == strings.ToUpper(nombre) {
			return d.Distance
		}
	}
	return -1

}

//Procesar la informacion de los satelites
func Procesar() (resp responseDecoded, err int) {
	var distancias []float32
	var mensajes [][]string
	err = 0
	//var respuesta responseDecoded
	if len(sat.Satellite) == 3 {
		dist := Find(sat.Satellite, "Kenobi")
		if dist >= 0 {
			distancias = append(distancias, dist)
		} else {
			err = 1
		}
		dist = Find(sat.Satellite, "Skywalker")
		if dist >= 0 {
			distancias = append(distancias, Find(sat.Satellite, "Skywalker"))
		} else {
			err = 1
		}
		dist = Find(sat.Satellite, "Sato")
		if dist >= 0 {
			distancias = append(distancias, Find(sat.Satellite, "Sato"))
		} else {
			err = 1
		}
	} else {
		err = 1
	}
	for _, d := range distancias {
		if d < 0 {
			err = 1
		}
	}
	if err == 0 {
		for _, d := range sat.Satellite {
			//distancias = append(distancias, d.Distance)
			mensajes = append(mensajes, d.Message)
		}
		x, y := GetLocation(distancias...)

		if x == 100 && y == -100 {
			err = 1
		}

		resp.Pos = position{X: x, Y: y}
		resp.Message = GetMessage(mensajes...)

	} else {
		err = 1
	}

	return resp, err
}
