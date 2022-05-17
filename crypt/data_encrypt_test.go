package crypt

import (
	"sort"
	"strings"
	"testing"
)

const (
	fishText1 = `Ясность нашей позиции очевидна: постоянное информационно-пропагандистское обеспечение нашей деятельности создаёт необходимость включения в производственный план целого ряда внеочередных мероприятий с учётом комплекса распределения внутренних резервов и ресурсов. Приятно, граждане, наблюдать, как стремящиеся вытеснить традиционное производство, нанотехнологии представляют собой не что иное, как квинтэссенцию победы маркетинга над разумом и должны быть своевременно верифицированы. Лишь ключевые особенности структуры проекта, инициированные исключительно синтетически, своевременно верифицированы! Как принято считать, активно развивающиеся страны третьего мира, инициированные исключительно синтетически, разоблачены.`
	fishText2 = `Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a, tellus. Phasellus viverra nulla ut metus varius laoreet. Quisque rutrum. Aenean imperdiet. Etiam ultricies nisi vel augue. Curabitur ullamcorper ultricies nisi. Nam eget dui. Etiam rhoncus. Maecenas tempus, tellus eget condimentum rhoncus, sem quam semper libero, sit amet adipiscing sem neque sed ipsum. Nam quam nunc, blandit vel, luctus pulvinar, hendrerit id, lorem. Maecenas nec odio et ante tincidunt tempus. Donec vitae sapien ut libero venenatis faucibus. Nullam quis ante. Etiam sit amet orci eget eros faucibus tincidunt. Duis leo. Sed fringilla mauris sit amet nibh. Donec sodales sagittis magna. Sed consequat, leo eget bibendum sodales, augue velit cursus nunc.`
	fishText3 = `Далеко-далеко за словесными горами в стране гласных и согласных живут рыбные тексты. Вдали от всех живут они в буквенных домах на берегу Семантика большого языкового океана.`
	fishText4 = `Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem. Nulla consequat massa quis enim. Donec pede justo, fringilla vel, aliquet nec, vulputate eget, arcu. In enim justo, rhoncus ut, imperdiet a, venenatis vitae, justo. Nullam dictum felis eu pede mollis pretium. Integer tincidunt. Cras dapibus. Vivamus elementum semper nisi. Aenean vulputate eleifend tellus. Aenean leo ligula, porttitor eu, consequat vitae, eleifend ac, enim. Aliquam lorem ante, dapibus in, viverra quis, feugiat a`
)

func TestEncryptDecryptData(t *testing.T) {
	t.Parallel()
	type args struct {
		data string
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		entV    int
		want    string
		wantErr bool
	}{
		{
			name: "test short data",
			args: args{
				data: "some word for encoding",
				key:  []byte("abcdefghijklmnopqrstuvwxyz"),
			},
			entV:    42,
			want:    "some word for encoding",
			wantErr: false,
		},
		{
			name: "test fish text 1",
			args: args{
				data: fishText1,
				key:  []byte(fishText2),
			},
			entV:    8,
			want:    fishText1,
			wantErr: false,
		},
		{
			name: "test fish text 2",
			args: args{
				data: fishText3,
				key:  []byte(fishText4),
			},
			entV:    8,
			want:    fishText3,
			wantErr: false,
		},
		{
			name: "test fish text with simple key",
			args: args{
				data: fishText3,
				key:  []byte(strings.Repeat("abc", len(fishText3))),
			},
			entV:    8,
			want:    fishText3,
			wantErr: false,
		},
		{
			name: "short key error",
			args: args{
				data: "some text for encoding. this text is longer than its key. expected error",
				key:  []byte("some short key"),
			},
			wantErr: true,
		},
	}
	for i := range tests {
		var test = tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var data = []byte(test.args.data)
			err := EncryptDecryptData(data, test.args.key)
			if (err != nil) != test.wantErr {
				t.Errorf("EncryptDecryptData() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if err != nil {
				return
			}
			if strings.EqualFold(test.want, string(data)) {
				t.Errorf("the data is not encoded: %s", test.want)
			}
			if ent := checkEncodedEntropy(data); ent > test.entV {
				t.Errorf("insufficient entropy value: %d", ent)
			}
			err = EncryptDecryptData(data, test.args.key)
			if err != nil {
				t.Errorf("EncryptDecryptData() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !strings.EqualFold(test.want, string(data)) {
				t.Errorf("the data is not decoded: %s", test.want)
			}
		})
	}
}

func checkEncodedEntropy(data []byte) (result int) {
	var max, min byte = 0, 255
	var ent = make(map[byte]int, 0)
	for _, b := range data {
		ent[b]++
		if max < b {
			max = b
		}
		if min > b {
			min = b
		}
	}
	var keys []int
	for k := range ent {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	var l = 0
	for _, k := range keys {
		val := ent[byte(k)] - l
		l = ent[byte(k)]
		result += val
	}
	return result * (256 - int(max-min))
}
