package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const contentFour = "cryptocurrency!"
const contentFive = "cryptocurrency!"

// const contentFive = "cryptocurrency!"
// const contentSix = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var doris *client.User
	var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User
	var charlesDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("", defaultPassword)
			Expect(err).NotTo(BeNil())

			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			alice, err = client.InitUser("alice", "g")
			Expect(err).NotTo(BeNil())

			aliceLaptop, err = client.GetUser("alice", "fff")
			Expect(err).NotTo(BeNil())

			aliceLaptop, err = client.GetUser("alcice", "fff")
			Expect(err).NotTo(BeNil())

			err = alice.StoreFile(aliceFile, []byte(""))
			Expect(err).To(BeNil())

		})

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("dsadsadsadsadasdasdasdsdsadsadsadsadsadsasdsadsasadsdsadasdsadsadsadsasdasdasdsadsadsadasddassdasdasdsasaddassdsadsadasdasdasdsadsadsad", defaultPassword)
			Expect(err).To(BeNil())

			aliceLaptop, err = client.GetUser("dsadsadsadsadasdasdasdsdsadsadsadsadsadsasdsadsasadsdsadasdsadsadsadsasdasdasdsadsadsadasddassdasdasdsasaddassdsadsadasdasdasdsadsadsad", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile("ssadasdsadasdassadasdsadasdasdsdsadassadsadasdasdasdasdasdasdassadsadasdasdsadasdasdasdasdasdasdsadasdadasdasdasdaasdasdadsdasdadadasdasdasasdsaasd", []byte(contentOne))
			Expect(err).To(BeNil())

			data, err := alice.LoadFile("ssadasdsadasdassadasdsadasdasdsdsadassadsadasdasdasdasdasdasdassadsadasdasdsadasdasdasdasdasdasdsadasdadasdasdasdaasdasdadsdasdadadasdasdasasdsaasd")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")

			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", "")
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).NotTo(BeNil())

			data, err := charles.LoadFile(charlesFile)
			Expect(err).NotTo(BeNil())
			Expect(data).NotTo(Equal([]byte(contentOne)))

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).NotTo(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			err = bob.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())
		})
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("", defaultPassword)
			Expect(err).NotTo(BeNil())

			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", "")
			Expect(err).To(BeNil())

			charlesDesktop, err = client.GetUser("charles", "")
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			invite, err = alice.CreateInvitation(aliceFile, "doris")
			Expect(err).To(BeNil())

			err = doris.AcceptInvitation("alice", invite, dorisFile)
			Expect(err).To(BeNil())

			//check
			data, err := doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			err = bob.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentFour))
			Expect(err).To(BeNil())

			err = doris.AppendToFile(dorisFile, []byte(contentFive))
			Expect(err).To(BeNil())

			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			//revoke
			err = bob.RevokeAccess(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(aliceFile, "boob")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
			Expect(data).ToNot(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())
			Expect(data).ToNot(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = charlesDesktop.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())
			Expect(data).ToNot(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			// add more
			aliceLaptop, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).NotTo(BeNil())

			err = bob.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).NotTo(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).NotTo(BeNil())
			//

			eve, err = client.InitUser("eve", defaultPassword)
			Expect(err).To(BeNil())

			invite, err = doris.CreateInvitation(dorisFile, "eve")
			Expect(err).To(BeNil())

			err = eve.AcceptInvitation("doris", invite, "hi")
			Expect(err).To(BeNil())

			data, err = eve.LoadFile("hi")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(aliceFile, "doris")
			Expect(err).To(BeNil())

			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = eve.LoadFile(dorisFile)
			Expect(err).NotTo(BeNil())
			Expect(data).NotTo(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))

			data, err = doris.LoadFile(dorisFile)
			Expect(err).NotTo(BeNil())
			Expect(data).NotTo(Equal([]byte(contentOne + contentTwo + contentThree + contentFour + contentFive)))
		})

		Specify("Basic Test: Multi level sharing.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			data11, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentOne)))

			err = aliceLaptop.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data11, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentOne + contentTwo)))

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			err = bob.StoreFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())

			invite, err := bob.CreateInvitation(bobFile, "alice")
			Expect(err).To(BeNil())

			err = alice.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			data11, err = alice.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentThree)))

			data11, err = aliceLaptop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentThree)))
		})

		Specify("Basic Test: Multi level sharing.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			err = aliceLaptop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			data11, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			err = bob.StoreFile(bobFile, []byte(contentFour))
			Expect(err).To(BeNil())

			invite, err := bob.CreateInvitation(bobFile, "alice")
			Expect(err).To(BeNil())

			err = alice.AcceptInvitation("bob", invite, aliceFile)
			Expect(err).NotTo(BeNil())

			err = alice.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			data11, err = aliceLaptop.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour)))

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			invite, err = alice.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			err = doris.AcceptInvitation("alice", invite, dorisFile)
			Expect(err).To(BeNil())

			data11, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour)))

			err = doris.AppendToFile(dorisFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data11, err = doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour + contentTwo)))

			data11, err = alice.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour + contentTwo)))

			data11, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour + contentTwo)))

			err = alice.RevokeAccess(charlesFile, "doris")
			Expect(err).NotTo(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			data11, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data11).To(Equal([]byte(contentFour + contentTwo)))

			err = bob.RevokeAccess(bobFile, "alice")
			Expect(err).To(BeNil())
		})
		// inituser, getuser
		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
		})
		// appendtofile, loadfile when not received

		//check -> appendtofile, loadfile when received, createinvitation, receivedinvitation
		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			aliceLaptop, err = client.InitUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

			alicePhone, err = client.InitUser("", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			invite2, err := bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite2, charlesFile)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			invite3, err := aliceLaptop.CreateInvitation(aliceFile, "doris")
			Expect(err).To(BeNil())

			err = doris.AcceptInvitation("alice", invite3, dorisFile)
			Expect(err).To(BeNil())

			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			data2, err := charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data2).To(Equal([]byte(contentOne)))

			data3, err := doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data3).To(Equal([]byte(contentOne)))

			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			data4, err := bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
			Expect(data4).ToNot(Equal([]byte(contentOne)))

			data5, err := charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())
			Expect(data5).ToNot(Equal([]byte(contentOne)))

			data6, err := doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = doris.AppendToFile(dorisFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data7, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data7).To(Equal([]byte(contentOne + contentTwo)))

			data8, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data8).To(Equal([]byte(contentOne + contentTwo)))

			err = aliceDesktop.StoreFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			data9, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data9).To(Equal([]byte(contentThree)))

			data10, err := doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data10).To(Equal([]byte(contentOne + contentTwo)))
		})

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("Alie", "")
			Expect(err).ToNot(BeNil())
		})
		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation("fff", "bob")
			Expect(err).NotTo(BeNil())

			invite, err = alice.CreateInvitation(aliceFile, "ddbob")
			Expect(err).NotTo(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).NotTo(BeNil())

			invite, err = alice.CreateInvitation("fff", "bdob")
			Expect(err).NotTo(BeNil())

			err = bob.AcceptInvitation("aliccce", invite, bobFile)
			Expect(err).NotTo(BeNil())

			err = bob.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err = bob.CreateInvitation(bobFile, "alice")
			Expect(err).To(BeNil())

			err = alice.AcceptInvitation("bob", invite, aliceFile)
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(aliceFile, "doris")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).NotTo(BeNil())

			err = alice.AcceptInvitation("bob", invite, bobFile)
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(bobFile, "doris")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).NotTo(BeNil())

			invite, err = alice.CreateInvitation(bobFile, "cf")
			Expect(err).NotTo(BeNil())

			invite, err = alice.CreateInvitation(charlesFile, "bob")
			Expect(err).NotTo(BeNil())

			invite, err = alice.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).To(BeNil())

			err = alice.RevokeAccess(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).NotTo(BeNil())

			err = charles.RevokeAccess(charlesFile, "bob")
			Expect(err).NotTo(BeNil())

			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = bob.RevokeAccess(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			err = bob.RevokeAccess(bobFile, "alice")
			Expect(err).To(BeNil())

			invite, err = alice.CreateInvitation(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			invite, err = aliceLaptop.CreateInvitation(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			invite, err = aliceLaptop.CreateInvitation(bobFile, "charles")
			Expect(err).NotTo(BeNil())

			err = aliceLaptop.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).NotTo(BeNil())

			err = aliceLaptop.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err = aliceLaptop.CreateInvitation(bobFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, charlesFile)
			Expect(err).To(BeNil())

			err = aliceLaptop.RevokeAccess(bobFile, "bob")
			Expect(err).To(BeNil())

		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			data2, err := alice.LoadFile(charlesFile)
			Expect(err).NotTo(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			invite2, err := bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite2, charlesFile)
			Expect(err).To(BeNil())

			doris, err = client.InitUser("doris", defaultPassword)
			Expect(err).To(BeNil())

			invite3, err := alice.CreateInvitation(aliceFile, "doris")
			Expect(err).To(BeNil())

			err = doris.AcceptInvitation("alice", invite3, dorisFile)
			Expect(err).To(BeNil())

			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			data2, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data2).To(Equal([]byte(contentOne)))

			data3, err := doris.LoadFile(dorisFile)
			Expect(err).To(BeNil())
			Expect(data3).To(Equal([]byte(contentOne)))

			err = alice.RevokeAccess(aliceFile, "doris")
			Expect(err).To(BeNil())

			data4, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data4).To(Equal([]byte(contentOne)))

			data5, err := charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data5).To(Equal([]byte(contentOne)))

			data6, err := doris.LoadFile(dorisFile)
			Expect(err).ToNot(BeNil())
			Expect(data6).ToNot(Equal([]byte(contentOne)))
		})

		// Init User

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("Alice", "")
			Expect(err).ToNot(BeNil())

			aliceLaptop, err = client.GetUser("", "")
			Expect(err).ToNot(BeNil())

			aliceLaptop, err = client.GetUser("", defaultPassword)
			Expect(err).ToNot(BeNil())

			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		//getuser
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			aliceDesktop, err = client.GetUser("alice", "wrong")
			Expect(err).ToNot(BeNil())

			aliceDesktop, err = client.GetUser("alice", "")
			Expect(err).ToNot(BeNil())

			aliceDesktop, err = client.GetUser("Alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		//store
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			err = aliceLaptop.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data6, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentTwo)))

			data6, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentTwo)))

			err = alice.StoreFile("3", []byte(contentOne))
			Expect(err).To(BeNil())

			data6, err = alice.LoadFile("3")
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			// err = alice.StoreFile("", []byte(contentOne))
			// Expect(err).To(BeNil())

			// data6, err = alice.LoadFile("")
			// Expect(err).To(BeNil())
			// Expect(data6).To(Equal([]byte(contentOne)))
		})
		// load, append
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			data6, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			data6, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			data6, err = aliceLaptop.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
			Expect(data6).ToNot(Equal([]byte(contentOne)))

			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data6, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne + contentTwo)))

			data6, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne + contentTwo)))

			err = aliceLaptop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			data6, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			data6, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))
		})
		// invitation

		//multi level
		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			data6, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			data6, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne)))

			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			data6, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne + contentTwo)))

			data6, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne + contentTwo)))

			data6, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data6).To(Equal([]byte(contentOne + contentTwo)))
		})
	})
})
