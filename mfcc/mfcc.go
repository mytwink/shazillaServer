package mfcc

import(
	"math"
	"fmt"
	"log"
)

const (
	FRAMESIZE = 200
	FILTERCOUNT = 12
	FLSTART = 300
	FLEND = 8000
	COEFCOUNT = 10
)

type MySample struct {
	Id int64
	Name,Path string
	Vector []float64
	Dist float64
}

type MySimpleSample struct {
	Id int64
	Name,Path string
	Dist float64
}

type Mfcc struct {
	FileSize,SampleRate int64
	File,Re,Im,Fsmp []int64
	Pn,ScaleCenter,Xn []float64
	Cn,Hik [][]float64
}

func NewMfcc(data []int64, sampleRate int64) *Mfcc {
	mfcc := &Mfcc{}
	mfcc.File = data
	mfcc.FileSize = int64(len(data))
	mfcc.SampleRate = sampleRate

	var i int64

	mfcc.SetCenters()
	log.Println("SetCenters")
	mfcc.SetSmp()
	log.Println("SetSmp")
	mfcc.SetHik()
	log.Println("SetHik")
	for i=0; i<mfcc.FileSize; i=i+int64(FRAMESIZE/2) {
		if(i+int64(FRAMESIZE) > mfcc.FileSize) {
			break;
		}

		frame := []int64{}

		for j:=0; j<FRAMESIZE; j++ {
			frame = append(frame, mfcc.File[int(i)+j])
		}

		mfcc.Pn = []float64{}

		for k:=0; k<FRAMESIZE; k++ {
			mfcc.Pn = append(mfcc.Pn, Fourier(frame,int64(k)))
		}

		mfcc.SetXn()
		mfcc.getCepstral()
	}

	mfcc.File = []int64{}
	mfcc.Pn = []float64{}
	mfcc.ScaleCenter = []float64{}
	mfcc.Fsmp = []int64{}
	mfcc.Xn = []float64{}
	mfcc.Hik = [][]float64{}

	return mfcc
}

func Fourier(input []int64,k int64) float64 {
	var re,im float64
	for i:=0; i<FRAMESIZE; i++ {
		re+=float64(input[i])*math.Cos(math.Pi*float64(2*int(k)*i)/float64(FRAMESIZE))
		im+=float64(input[i])*math.Sin(math.Pi*float64(2*int(k)*i)/float64(FRAMESIZE))
	}

	re = (2.0/float64(FRAMESIZE))*re
	im = -(2.0/float64(FRAMESIZE))*im
	return re*re+im*im
}

func Mel2Herz(mel float64) float64 {
	return 700.0*(math.Exp(mel/1127.0)-1)
}

func Herz2Mel(herz float64) float64 {
	return 1127*math.Log(1+herz/700.0)
}

func (mfcc *Mfcc) SetCenters() {
	mfcc.ScaleCenter = make([]float64, FILTERCOUNT)
	fl_l := Herz2Mel(FLSTART)
	fl_h := Herz2Mel(FLEND)
	l := (fl_h-fl_l) / float64(FILTERCOUNT-1)

	for i:=0; i<FILTERCOUNT; i++ {
		mfcc.ScaleCenter[i] = Mel2Herz(fl_l+l*float64(i))
	}
}

func (mfcc *Mfcc) SetSmp() {
	mfcc.Fsmp = make([]int64, FILTERCOUNT)
	for i:=0; i<FILTERCOUNT; i++ {
		mfcc.Fsmp[i] = int64(math.Floor(float64(FRAMESIZE+1)*mfcc.ScaleCenter[i]/float64(mfcc.SampleRate)))
	}
}

func (mfcc *Mfcc) SetHik() {
	mfcc.Hik = make([][]float64, COEFCOUNT)
	var val float64
	for i:=1; i<=COEFCOUNT; i++ {
		temp := make([]float64, FRAMESIZE)

		var k int64
		for k=0; k<FRAMESIZE; k++ {
			if k < mfcc.Fsmp[i-1] {
				val = 0.0
			}

			if k >= mfcc.Fsmp[i-1] && k <= mfcc.Fsmp[i] {
				val = float64(k-mfcc.Fsmp[i-1])/float64(mfcc.Fsmp[i]-mfcc.Fsmp[i-1])
			}

			if k <= mfcc.Fsmp[i+1] && k >= mfcc.Fsmp[i] {
				val = float64(mfcc.Fsmp[i+1]-k)/float64(mfcc.Fsmp[i+1]-mfcc.Fsmp[i])
			}

			if k > mfcc.Fsmp[i+1] {
				val = 0.0
			}

			temp[k] = val
		}

		mfcc.Hik[i-1] =	temp
	}
}

func (mfcc *Mfcc) SetXn() {
	var summ float64

	mfcc.Xn = make([]float64, COEFCOUNT)

	for i:=1; i<=COEFCOUNT; i++ {
		summ = 0.0
		for k:=0; k<FRAMESIZE; k++ {
			summ+=mfcc.Pn[k]*mfcc.Hik[i-1][k]
		}
		mfcc.Xn[i-1] = math.Log(summ)
	}
}

func (mfcc *Mfcc) getCepstral() {
	summ := 0.0
	temp := make([]float64, COEFCOUNT)

	for j:=0; j<COEFCOUNT; j++ {
		summ = 0.0
		for k:=0; k < COEFCOUNT; k++ {
			summ += mfcc.Xn[k]*math.Cos(float64(j*(2*k+1))*math.Pi/float64(COEFCOUNT*2))
		}
		temp[j] = summ
	}

	mfcc.Cn = append(mfcc.Cn, temp)
}

func getDis(vec1 []float64, offset1 int, vec2 []float64, offset2, len int) float64 {
	var dis float64
	for i:=0; i<len; i++ {
		dis += (vec1[offset1+i] - vec2[offset2+i]) * (vec1[offset1+i] - vec2[offset2+i])
	}
	return math.Sqrt(dis)
}

func (mfcc *Mfcc) GetVector() []float64 {
	vector := make([]float64, len(mfcc.Cn)*COEFCOUNT)
	for i, vec := range mfcc.Cn {
		for j, elem := range vec {
			vector[i*COEFCOUNT+j] = elem
		}
	}
	return vector
}

func (mfcc *Mfcc) Chisqr(SubstractArr []float64, sample MySimpleSample, c chan<- MySimpleSample) {
	vec := mfcc.GetVector()

	var dis float64

	dim1 := len(vec)/COEFCOUNT
	dim2 := len(SubstractArr)/COEFCOUNT

	dp := make([][]float64, 2)
	dp[0] = make([]float64, dim2)
	dp[1] = make([]float64, dim2)

	for j:=0; j<dim2; j++ {
		dp[0][j] = getDis(vec, 0, SubstractArr, j*COEFCOUNT, COEFCOUNT)
	}

	for i:=1; i<dim1; i++ {
		for j:=0; j<dim2; j++ {
			dis = getDis(vec, i*COEFCOUNT, SubstractArr, j*COEFCOUNT, COEFCOUNT)

			dp[1][j] = dp[0][j]

			if j>0 {
				if dp[0][j-1] < dp[1][j] {
					dp[1][j] = dp[0][j-1]
				}
				if dp[1][j-1] < dp[1][j] {
					dp[1][j] = dp[1][j-1]
				}
			}
			dp[1][j] += dis
		}
		for j:=0; j<dim2; j++ {
			dp[0][j] = dp[1][j]
		}
	}
	dis = dp[0][dim2-1]
	sample.Dist = dis

	c <- sample
}

type emptyArrayError struct {
	msg string
}

func (e emptyArrayError) Error() string {
	return fmt.Sprintf("%v", e.msg)
}