package main

import "testing"

func TestFormatText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "basic",
			text:     "M=0\n",
			expected: "M=0",
		},
		{
			name:     "space",
			text:     " M = D + 1 \n",
			expected: "M=D+1",
		},
		{
			name:     "newline",
			text:     "\n",
			expected: "",
		},
		{
			name:     "commend",
			text:     "M=0 // comment",
			expected: "M=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := formatText(tt.text)

			if command != tt.expected {
				t.Errorf("expected %s but got %s\n", tt.expected, command)
			}
		})
	}
}

func TestGetCommandType(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected int
	}{
		{
			name:     "A command",
			command:  "@100",
			expected: A_COMMAND,
		},
		{
			name:     "C command",
			command:  "M=0",
			expected: C_COMMAND,
		},
		{
			name:     "symbol",
			command:  "(LOOP)",
			expected: SYMBOL,
		},
		{
			name:     "invalid command type",
			command:  "kontya",
			expected: INVALID_COMMAND,
		},
		{
			name:     "A command use variable",
			command:  "@i",
			expected: A_COMMAND,
		},
		{
			name:     "C command use operation",
			command:  "M=D+1",
			expected: C_COMMAND,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandType := getCommandType(tt.command)

			if commandType != tt.expected {
				t.Errorf("expected %d but got %d\n", tt.expected, commandType)
			}
		})
	}
}

func TestGetLabelSymbol(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "basic",
			command:  "(LOOP)",
			expected: "LOOP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbol := getLabelSymbol(tt.command)

			if symbol != tt.expected {
				t.Errorf("expected %s but got %s\n", tt.expected, symbol)
			}
		})
	}
}

func TestGetSymbolAddress(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		expected int
	}{
		{
			name:     "basic",
			symbol:   "LOOP",
			expected: 16,
		},
		{
			name:     "2nd symbol",
			symbol:   "HOGE",
			expected: 17,
		},
		{
			name:     "symbol exist",
			symbol:   "LOOP",
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := getSymbolAddress(tt.symbol)

			if address != tt.expected {
				t.Errorf("expected %d but got %d\n", tt.expected, address)
			}
		})
	}
}

func TestGetValue_OK(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected int
	}{
		{
			name:     "basic",
			command:  "@100",
			expected: 100,
		},
		{
			name:     "symbol",
			command:  "@i",
			expected: nowSymbolAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := getValue(tt.command)
			if value != tt.expected {
				t.Errorf("expected %d but got %d", tt.expected, value)
			}
		})
	}
}

func TestFormatValue_OK(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "basic",
			value:    100,
			expected: "000000001100100",
		},
		{
			name:     "max value",
			value:    32767,
			expected: "111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valueString, err := formatValue(tt.value)
			if err != nil {
				t.Errorf("expected %s but got %v", tt.expected, err)
			}

			if valueString != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, valueString)
			}
		})
	}
}

func TestFormatValue_NG(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "over flow",
			value:    32768,
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valueString, err := formatValue(tt.value)
			if err == nil {
				t.Errorf("expected %s but got %s", tt.expected, valueString)
			}
		})
	}
}

func TestGetMnemonic_OK(t *testing.T) {
	type Expected struct {
		comp string
		dest string
		jump string
	}

	tests := []struct {
		name     string
		command  string
		expected Expected
	}{
		{
			name:    "operation command",
			command: "M=1",
			expected: Expected{
				comp: "1",
				dest: "M",
				jump: "",
			},
		},
		{
			name:    "jump command",
			command: "0;JMP",
			expected: Expected{
				comp: "0",
				dest: "",
				jump: "JMP",
			},
		},
		{
			name:    "complex operation command",
			command: "M=D+1",
			expected: Expected{
				comp: "D+1",
				dest: "M",
				jump: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, dest, jump, err := getMnemonic(tt.command)
			if err != nil {
				t.Errorf("expected %s, %s, %s but got %v", tt.expected.comp, tt.expected.dest, tt.expected.jump, err)
			}

			if comp != tt.expected.comp {
				t.Errorf("expected %s but got %s\n", tt.expected.comp, comp)
			}
			if dest != tt.expected.dest {
				t.Errorf("expected %s but got %s\n", tt.expected.dest, dest)
			}
			if jump != tt.expected.jump {
				t.Errorf("expected %s but got %s\n", tt.expected.jump, jump)
			}
		})
	}
}

func TestGetMnemonic_NG(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "invalid command",
			command:  "kontya",
			expected: "error",
		},
		{
			name:     "invalid jump command",
			command:  ";JMP",
			expected: "error",
		},
		{
			name:     "invalid jump command",
			command:  "0;",
			expected: "error",
		},
		{
			name:     "invalid operation command",
			command:  "M=",
			expected: "error",
		},
		{
			name:     "invalid operation command",
			command:  "=1",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, dest, jump, err := getMnemonic(tt.command)
			if err == nil {
				t.Errorf("expected %s but got %s, %s, %s", tt.expected, comp, dest, jump)
			}
		})
	}
}

func TestGetComp_OK(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "0",
			mnemonic: "0",
			expected: "0101010",
		},
		{
			name:     "1",
			mnemonic: "1",
			expected: "0111111",
		},
		{
			name:     "-1",
			mnemonic: "-1",
			expected: "0111010",
		},
		{
			name:     "D",
			mnemonic: "D",
			expected: "0001100",
		},
		{
			name:     "A",
			mnemonic: "A",
			expected: "0110000",
		},
		{
			name:     "!D",
			mnemonic: "!D",
			expected: "0001101",
		},
		{
			name:     "!A",
			mnemonic: "!A",
			expected: "0110001",
		},
		{
			name:     "-D",
			mnemonic: "-D",
			expected: "0001111",
		},
		{
			name:     "-A",
			mnemonic: "-A",
			expected: "0110011",
		},
		{
			name:     "D+1",
			mnemonic: "D+1",
			expected: "0011111",
		},
		{
			name:     "A+1",
			mnemonic: "A+1",
			expected: "0110111",
		},
		{
			name:     "D-1",
			mnemonic: "D-1",
			expected: "0001110",
		},
		{
			name:     "A-1",
			mnemonic: "A-1",
			expected: "0110010",
		},
		{
			name:     "D+A",
			mnemonic: "D+A",
			expected: "0000010",
		},
		{
			name:     "D-A",
			mnemonic: "D-A",
			expected: "0010011",
		},
		{
			name:     "A-D",
			mnemonic: "A-D",
			expected: "0000111",
		},
		{
			name:     "D&A",
			mnemonic: "D&A",
			expected: "0000000",
		},
		{
			name:     "D|A",
			mnemonic: "D|A",
			expected: "0010101",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := getComp(tt.mnemonic)
			if err != nil {
				t.Errorf("expected %s but got %v\n", tt.expected, err)
			}

			if comp != tt.expected {
				t.Errorf("expected %s but got %s\n", tt.expected, comp)
			}
		})
	}
}

func TestGetComp_NG(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "invalid mnemonic",
			mnemonic: "A+D",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp, err := getComp(tt.mnemonic)
			if err == nil {
				t.Errorf("expected %s but got %s\n", tt.expected, comp)
			}
		})
	}
}

func TestGetDest_OK(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "null",
			mnemonic: "",
			expected: "000",
		},
		{
			name:     "memory",
			mnemonic: "M",
			expected: "001",
		},
		{
			name:     "Dregister",
			mnemonic: "D",
			expected: "010",
		},
		{
			name:     "memory and Dregister",
			mnemonic: "MD",
			expected: "011",
		},
		{
			name:     "Aregister",
			mnemonic: "A",
			expected: "100",
		},
		{
			name:     "Aregister and memory",
			mnemonic: "AM",
			expected: "101",
		},
		{
			name:     "Aregister and Dregister",
			mnemonic: "AD",
			expected: "110",
		},
		{
			name:     "Aregister and memory and Dregister",
			mnemonic: "AMD",
			expected: "111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest, err := getDest(tt.mnemonic)
			if err != nil {
				t.Errorf("expected %s but got %v\n", tt.expected, err)
			}

			if dest != tt.expected {
				t.Errorf("expected %s but got %s\n", tt.expected, dest)
			}
		})
	}
}

func TestGotDest_NG(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "invalid mnemonic",
			mnemonic: "DM",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest, err := getDest(tt.mnemonic)
			if err == nil {
				t.Errorf("expected %s but got %s", tt.expected, dest)
			}
		})
	}
}

func TestGetJump_OK(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "null",
			mnemonic: "",
			expected: "000",
		},
		{
			name:     "more",
			mnemonic: "JGT",
			expected: "001",
		},
		{
			name:     "equal",
			mnemonic: "JEQ",
			expected: "010",
		},
		{
			name:     "more or equal",
			mnemonic: "JGE",
			expected: "011",
		},
		{
			name:     "less",
			mnemonic: "JLT",
			expected: "100",
		},
		{
			name:     "not equal",
			mnemonic: "JNE",
			expected: "101",
		},
		{
			name:     "less or equal",
			mnemonic: "JLE",
			expected: "110",
		},
		{
			name:     "anyway",
			mnemonic: "JMP",
			expected: "111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jump, err := getJump(tt.mnemonic)
			if err != nil {
				t.Errorf("expected %s but got %v\n", tt.expected, err)
			}

			if jump != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, jump)
			}
		})
	}
}

func TestGetJump_NG(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		expected string
	}{
		{
			name:     "invalid mnemonic",
			mnemonic: "JUMP",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jump, err := getJump(tt.mnemonic)
			if err == nil {
				t.Errorf("expected %s but got %s\n", tt.expected, jump)
			}
		})
	}
}
