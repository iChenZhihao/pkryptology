package node

import (
	"fmt"
	tt "github.com/coinbase/kryptology/internal"
	"github.com/coinbase/kryptology/pkg/core"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVerifySigNodes(t *testing.T) {
	PublicKey := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("74583413151673891188237912075727467634522453383318553670848436079828535975354"),
		Y:     tt.B10("110942599231583866555011384577907428229443523206202810845808044399708421596435"),
	}
	sigma := &curves.EcdsaSignature{
		V: 1,
		R: tt.B10("100126053408351854952837750211155726448882031017382074913016609328584383795848"),
		S: tt.B10("54835403759799446985808067844383347762118791931607875585726963116856918633647"),
	}
	bytes := []byte("Hello World!")
	b, _ := core.Hash(bytes, curve)
	verifier := ecdsaVerifier(PublicKey, b.Bytes(), sigma)
	fmt.Println(verifier)
	require.True(t, verifier)
}

func TestLocalRounds(t *testing.T) {
	PublicKey := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("39358771517832395540406228045245140109111041015819744873599548220653790899702"),
		Y:     tt.B10("47321279348055860705548442743973923714922108035385310788178691142457952163074"),
	}
	sigma := &curves.EcdsaSignature{
		V: 1,
		R: tt.B10("25909809917719420636663693830758021531901966845092433536299352386642173978062"),
		S: tt.B10("55038030789096116953041213828106061762776009669336107803119278065672478417446"),
	}
	bytes := []byte("Hello World!")
	b, _ := core.Hash(bytes, curve)
	verifier := ecdsaVerifier(PublicKey, b.Bytes(), sigma)
	fmt.Println(verifier)
	require.True(t, verifier)
}

func TestRbarkVerify(t *testing.T) {
	Rbark1 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("43366393737367283506230664002259962261368750878449170150625392662270619635419"),
		Y:     tt.B10("79222081361407867378038332722817677767497652584542000982598777430265810227090"),
	}

	Rbark2 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("92386699970857092351543176696544623895720642266657378525854638114838122959952"),
		Y:     tt.B10("4906279506807375338154670360048242739676136322961374288096888574611966053222"),
	}

	Rbark3 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("30007084835755579160566572104698697057736316134815105471320967591763329146340"),
		Y:     tt.B10("58731769056084816811333706700625351955766090497207182653414722364624200646529"),
	}

	pk := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("96955889283678400263068037461262761680914128219188159310548473526128827365218"),
		Y:     tt.B10("9085252391252934592332740421031848998887611779018991136705027281696554466876"),
	}

	add, err := Rbark1.Add(Rbark2)
	require.NoError(t, err)
	Rbark, err := add.Add(Rbark3)
	require.NoError(t, err)

	Rbark.Y, err = core.Neg(Rbark.Y, curve.Params().P)
	require.NoError(t, err)
	Rbark, err = Rbark.Add(pk)
	require.NoError(t, err)

	if !Rbark.IsIdentity() {
		t.Errorf("%v != %v", Rbark.X, pk.X)
		t.FailNow()
	}

}

func TestLocalRbarkVerify(t *testing.T) {
	Rbark1 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("104912165187826790899160376665611661758243192245706652422853306095655885993307"),
		Y:     tt.B10("65156954815844156942666159360622801056451345932378969688663819551104666343756"),
	}

	Rbark2 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("45876822809308947617709134959562272808754391646327594096273060191149589459499"),
		Y:     tt.B10("59325007059639563417014080075588749734212049158765589321787306204360650851107"),
	}

	Rbark3 := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("82435865851825732851507815470061419622536099812499303970573872572054076013210"),
		Y:     tt.B10("33933609595237516634587849255289297869667530386501329044903568804302203051541"),
	}

	pk := &curves.EcPoint{
		Curve: curve,
		X:     tt.B10("1407748129790977477106506368293085372447119570463534828980860867521365067756"),
		Y:     tt.B10("68128200143010691737411761915219204694834847074129004209615941147772933462071"),
	}

	add, err := Rbark1.Add(Rbark2)
	require.NoError(t, err)
	Rbark, err := add.Add(Rbark3)
	require.NoError(t, err)

	Rbark.Y, err = core.Neg(Rbark.Y, curve.Params().P)
	require.NoError(t, err)
	Rbark, err = Rbark.Add(pk)
	require.NoError(t, err)

	if !Rbark.IsIdentity() {
		t.Errorf("%v != %v", Rbark.X, pk.X)
		t.FailNow()
	}

}
