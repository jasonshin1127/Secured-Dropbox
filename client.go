package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"
	// Optional.
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string
	Password string
	Priv_key userlib.PrivateKeyType
	Sign_key userlib.PrivateKeyType
	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

type File struct {
	K1             []byte
	K2             []byte
	Recipient_UUID map[string]uuid.UUID
	Recipient_HMAC map[string]uuid.UUID
	Recipient_K1   map[string][]byte
	Recipient_k2   map[string][]byte
	File_UUID      uuid.UUID
	HMAC_UUID      uuid.UUID
	Received       bool
}

type Invitation struct {
	K1        []byte
	K2        []byte
	File_UUID uuid.UUID
	HMAC_UUID uuid.UUID
	Index     int
	Reject    bool
}

type Invitation_info struct {
	Bytes_UUID []byte
	Bytes_HMAC []byte
	Bytes_K1   []byte
	Bytes_K2   []byte
	Sign_UUID  []byte
	Sign_HMAC  []byte
	Sign_K1    []byte
	Sign_K2    []byte
}

type ciphertext_index struct {
	Content []byte
	Index   int
}

// NOTE: The following methods have toy (insecure!) implementations.

func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	if len(username) == 0 {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	_, ok := userlib.KeystoreGet(username)
	if ok {
		return &userdata, errors.New(strings.ToTitle("invalid username"))
	}

	userdata.Username = username
	userdata.Password = password

	pub_key, priv_key, err := userlib.PKEKeyGen()
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	sign_key, ver_key, err := userlib.DSKeyGen()
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	userdata.Priv_key = priv_key
	userdata.Sign_key = sign_key

	userlib.KeystoreSet(username, pub_key)
	userlib.KeystoreSet(username+"ver_key", ver_key)

	byte_username, err := json.Marshal(username)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	byte_password, err := json.Marshal(password)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	iv := userlib.RandomBytes(16)

	key := userlib.Argon2Key(byte_username, byte_password, 16)
	key2 := userlib.Argon2Key(byte_password, byte_username, 16)

	userstruct_ekey, err := userlib.HashKDF(key, []byte("encryption"))
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	userstruct_ekey = userstruct_ekey[:16]
	userstruct_hkey, err := userlib.HashKDF(key, []byte("hmac"))
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	userstruct_hkey = userstruct_hkey[:16]

	byte_userdata, err := json.Marshal(userdata)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	ciphertext := userlib.SymEnc(userstruct_ekey, iv, byte_userdata)
	HMAC_key, err := userlib.HMACEval(userstruct_hkey, ciphertext)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	uuid_userdata, err := uuid.FromBytes(userlib.Hash(key)[:16])
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	uuid_HMAC, err := uuid.FromBytes(userlib.Hash(key2)[:16])
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	userlib.DatastoreSet(uuid_userdata, ciphertext)
	userlib.DatastoreSet(uuid_HMAC, HMAC_key)
	// HMAC_key

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User

	if len(username) == 0 {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	_, ok := userlib.KeystoreGet(username)
	if !ok {
		return &userdata, errors.New(strings.ToTitle("invalid username"))
	}

	byte_username, err := json.Marshal(username)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	byte_password, err := json.Marshal(password)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	key := userlib.Argon2Key(byte_username, byte_password, 16)
	key2 := userlib.Argon2Key(byte_password, byte_username, 16)
	userstruct_ekey, err := userlib.HashKDF(key, []byte("encryption"))
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	userstruct_ekey = userstruct_ekey[:16]
	userstruct_hkey, err := userlib.HashKDF(key, []byte("hmac"))
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}
	userstruct_hkey = userstruct_hkey[:16]

	uuid_userdata, err := uuid.FromBytes(userlib.Hash(key)[:16])
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	uuid_HMAC, err := uuid.FromBytes(userlib.Hash(key2)[:16])
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	ciphertext, ok := userlib.DatastoreGet(uuid_userdata)
	if !ok {
		return &userdata, errors.New(strings.ToTitle("invalid username or password"))
	}

	ds_HMAC_key, ok := userlib.DatastoreGet(uuid_HMAC)
	if !ok {
		return &userdata, errors.New(strings.ToTitle("Datastore Tamper!"))
	}

	HMAC_key, err := userlib.HMACEval(userstruct_hkey, ciphertext)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Invalid"))
	}

	check := userlib.HMACEqual(ds_HMAC_key, HMAC_key)
	if !check {
		return &userdata, errors.New(strings.ToTitle("Datastore Tamper!"))
	}

	marshal_user := userlib.SymDec(userstruct_ekey, ciphertext)
	err = json.Unmarshal(marshal_user, &userdata)
	if err != nil {
		return &userdata, errors.New(strings.ToTitle("Unmarshall Error"))
	}

	userdataptr = &userdata

	return userdataptr, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	var filedata File

	K1 := userlib.RandomBytes(16)
	K2 := userlib.RandomBytes(16)
	File_uuid := uuid.New()
	HMAC_uuid := uuid.New()

	filedata.K1 = K1
	filedata.K2 = K2
	filedata.File_UUID = File_uuid
	filedata.HMAC_UUID = HMAC_uuid
	filedata.Received = false

	Recipient_UUID := make(map[string]uuid.UUID)
	Recipient_HMAC := make(map[string]uuid.UUID)
	Recipient_K1 := make(map[string][]byte)
	Recipient_K2 := make(map[string][]byte)

	filedata.Recipient_UUID = Recipient_UUID
	filedata.Recipient_HMAC = Recipient_HMAC
	filedata.Recipient_K1 = Recipient_K1
	filedata.Recipient_k2 = Recipient_K2

	var ciphertext_index ciphertext_index
	ciphertext_index.Content = content
	ciphertext_index.Index = 0

	contentBytes, err := json.Marshal(ciphertext_index)
	if err != nil {
		return errors.New(strings.ToTitle(string(filedata.K2) + string(K2)))
	}
	iv := userlib.RandomBytes(16)

	file_ciphertext := userlib.SymEnc(K1, iv, contentBytes)
	file_HMAC, err := userlib.HMACEval(K2, file_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	userlib.DatastoreSet(File_uuid, file_ciphertext)
	userlib.DatastoreSet(HMAC_uuid, file_HMAC)

	// file struct storage
	byte_filedata, err := json.Marshal(filedata)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	bytes_storage_content := []byte(filename + userdata.Username)
	storage_content_uuid, err := uuid.FromBytes(userlib.Hash(bytes_storage_content)[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	bytes_HMAC_uuid := []byte(userdata.Username + filename)
	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash(bytes_HMAC_uuid)[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	storage_ciphertext := userlib.SymEnc(filestruct_ekey, iv, byte_filedata)
	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, storage_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	userlib.DatastoreSet(storage_content_uuid, storage_ciphertext)
	userlib.DatastoreSet(storage_HMAC_uuid, storage_HMAC)

	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) error {

	bytes_storage_content := []byte(filename + userdata.Username)
	storage_content_uuid, err := uuid.FromBytes(userlib.Hash(bytes_storage_content)[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	bytes_HMAC_uuid := []byte(userdata.Username + filename)
	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash(bytes_HMAC_uuid)[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	// checking HMAC
	ciphertext, ok := userlib.DatastoreGet(storage_content_uuid)
	if !ok {
		return errors.New(strings.ToTitle("1"))
	}

	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	ds_HMAC_key, ok := userlib.DatastoreGet(storage_HMAC_uuid)
	if !ok {
		return errors.New(strings.ToTitle("2"))
	}

	check := userlib.HMACEqual(ds_HMAC_key, storage_HMAC)
	if !check {
		return errors.New(strings.ToTitle("3"))
	}

	//Decrypting ciphertext
	var file File
	marshal_file := userlib.SymDec(filestruct_ekey, ciphertext)
	err = json.Unmarshal(marshal_file, &file)
	if err != nil {
		return errors.New(strings.ToTitle("Unmarshall Error"))
	}
	iv := userlib.RandomBytes(16)

	file_UUID := file.File_UUID
	HMAC_UUID := file.HMAC_UUID
	k1 := file.K1
	k2 := file.K2

	if file.Received {
		invitation_UUID := file.File_UUID
		invitation_HMAC := file.HMAC_UUID
		invitation_k1 := file.K1
		invitation_k2 := file.K2

		// invitation description
		invitation_ciphertext, ok := userlib.DatastoreGet(invitation_UUID)
		if !ok {
			return errors.New(strings.ToTitle("1"))
		}

		storage_invitation_HMAC, err := userlib.HMACEval(invitation_k2, invitation_ciphertext)
		if err != nil {
			return errors.New(strings.ToTitle("Invalid"))
		}

		ds_HMAC, ok := userlib.DatastoreGet(invitation_HMAC)
		if !ok {
			return errors.New(strings.ToTitle("2"))
		}

		check := userlib.HMACEqual(ds_HMAC, storage_invitation_HMAC)
		if !check {
			return errors.New(strings.ToTitle("3"))
		}

		//Decrypting ciphertext
		var invitation Invitation
		marshal_file := userlib.SymDec(invitation_k1, invitation_ciphertext)
		err = json.Unmarshal(marshal_file, &invitation)
		if err != nil {
			return errors.New(strings.ToTitle("Unmarshall Error"))
		}
		if invitation.Reject {
			return errors.New(strings.ToTitle("Revoked"))
		}

		file_UUID = invitation.File_UUID
		HMAC_UUID = invitation.HMAC_UUID
		k1 = invitation.K1
		k2 = invitation.K2
	}
	var original_ciphertext ciphertext_index
	//Retrieve actual file info
	file_ciphertext, ok := userlib.DatastoreGet(file_UUID)
	if !ok {
		return errors.New(strings.ToTitle(string(marshal_file)))
	}

	file_HMAC, err := userlib.HMACEval(k2, file_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("113"))
	}

	file_HMAC_key, ok := userlib.DatastoreGet(HMAC_UUID)
	if !ok {
		return errors.New(strings.ToTitle("11"))
	}

	check_file := userlib.HMACEqual(file_HMAC_key, file_HMAC)
	if !check_file {
		return errors.New(strings.ToTitle("Original File TTTampered"))
	}

	file_marshal := userlib.SymDec(k1, file_ciphertext)
	err = json.Unmarshal(file_marshal, &original_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Unmarshall Error"))
	}

	original_ciphertext.Index = original_ciphertext.Index + 1
	bytes_index, err := json.Marshal(original_ciphertext.Index)

	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	new_file_uuid, err := uuid.FromBytes(userlib.Hash([]byte(string(k1) + string(k2) + string(bytes_index)))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	new_HMAC_uuid, err := uuid.FromBytes(userlib.Hash([]byte(string(k2) + string(k1) + string(bytes_index)))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	var new_ciphertext ciphertext_index
	new_ciphertext.Content = content
	new_ciphertext.Index = original_ciphertext.Index

	bytes_new_ciphertext, err := json.Marshal(new_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	new_file_ciphertext := userlib.SymEnc(k1, iv, bytes_new_ciphertext)

	new_file_HMAC, err := userlib.HMACEval(k2, new_file_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	userlib.DatastoreSet(new_file_uuid, new_file_ciphertext)
	userlib.DatastoreSet(new_HMAC_uuid, new_file_HMAC)

	//restoring the updated ciphertext_index
	updated_contentBytes, err := json.Marshal(original_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("update fail"))
	}

	updated_file_ciphertext := userlib.SymEnc(k1, iv, updated_contentBytes)
	updated_file_HMAC, err := userlib.HMACEval(k2, updated_file_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	userlib.DatastoreSet(file_UUID, updated_file_ciphertext)
	userlib.DatastoreSet(HMAC_UUID, updated_file_HMAC)

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	contents := []byte{}

	bytes_storage_content := []byte(filename + userdata.Username)
	storage_content_uuid, err := uuid.FromBytes(userlib.Hash(bytes_storage_content)[:16])
	if err != nil {
		return contents, errors.New(strings.ToTitle("Invalid"))
	}

	bytes_HMAC_uuid := []byte(userdata.Username + filename)
	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash(bytes_HMAC_uuid)[:16])
	if err != nil {
		return contents, errors.New(strings.ToTitle("Invalid"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	// check HMAC
	ciphertext, ok := userlib.DatastoreGet(storage_content_uuid)
	if !ok {
		return contents, errors.New(strings.ToTitle("Original File Tampered1"))
	}

	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, ciphertext)
	if err != nil {
		return contents, errors.New(strings.ToTitle("44"))
	}

	ds_HMAC_key, ok := userlib.DatastoreGet(storage_HMAC_uuid)
	if !ok {
		return contents, errors.New(strings.ToTitle("Original File Tampered2"))
	}

	check := userlib.HMACEqual(ds_HMAC_key, storage_HMAC)
	if !check {
		return contents, errors.New(strings.ToTitle("Original File Tampered3"))
	}

	//Decrypting ciphertext of the file storage
	var file File
	marshal_file := userlib.SymDec(filestruct_ekey, ciphertext)
	err = json.Unmarshal(marshal_file, &file)
	if err != nil {
		return contents, errors.New(strings.ToTitle("Unmarshall Error"))
	}

	//file info
	k1 := file.K1
	k2 := file.K2
	file_UUID := file.File_UUID
	file_HMAC := file.HMAC_UUID

	//reading the file, going through the append list
	if file.Received {
		invitation_UUID := file.File_UUID
		invitation_HMAC := file.HMAC_UUID
		invitation_k1 := file.K1
		invitation_k2 := file.K2

		// invitation description
		invitation_ciphertext, ok := userlib.DatastoreGet(invitation_UUID)
		if !ok {
			return contents, errors.New(strings.ToTitle("1"))
		}

		storage_invitation_HMAC, err := userlib.HMACEval(invitation_k2, invitation_ciphertext)
		if err != nil {
			return contents, errors.New(strings.ToTitle("Invalid"))
		}

		ds_HMAC, ok := userlib.DatastoreGet(invitation_HMAC)
		if !ok {
			return contents, errors.New(strings.ToTitle("2"))
		}

		check := userlib.HMACEqual(ds_HMAC, storage_invitation_HMAC)
		if !check {
			return contents, errors.New(strings.ToTitle("3"))
		}

		//Decrypting ciphertext
		var invitation Invitation
		marshal_file := userlib.SymDec(invitation_k1, invitation_ciphertext)
		err = json.Unmarshal(marshal_file, &invitation)
		if err != nil {
			return contents, errors.New(strings.ToTitle("Unmarshall Error"))
		}

		if invitation.Reject {
			return contents, errors.New(strings.ToTitle("Revoked"))
		}
		file_UUID = invitation.File_UUID
		file_HMAC = invitation.HMAC_UUID
		k1 = invitation.K1
		k2 = invitation.K2
	}

	var content_index ciphertext_index
	//Retrieve actual file info
	file_ciphertext, ok := userlib.DatastoreGet(file_UUID)
	if !ok {
		return contents, errors.New(strings.ToTitle("Original File Tampered4"))
	}

	file_info_HMAC, err := userlib.HMACEval(k2, file_ciphertext)
	if err != nil {
		return contents, errors.New(strings.ToTitle("4"))
	}

	file_HMAC_key, ok := userlib.DatastoreGet(file_HMAC)
	if !ok {
		return contents, errors.New(strings.ToTitle("Original File Tampered5"))
	}

	check_file := userlib.HMACEqual(file_HMAC_key, file_info_HMAC)
	if !check_file {
		return contents, errors.New(strings.ToTitle("Original File Tampered6"))
	}

	file_marshal := userlib.SymDec(k1, file_ciphertext)
	err = json.Unmarshal(file_marshal, &content_index)
	if err != nil {
		return contents, errors.New(strings.ToTitle("Unmarshall Error"))
	}

	index := content_index.Index

	for i := 0; i <= index; i += 1 {
		if i == 0 {
			// HMAC
			file_ciphertext, ok := userlib.DatastoreGet(file_UUID)
			if !ok {
				return contents, errors.New(strings.ToTitle("5"))
			}

			created_file_HMAC, err := userlib.HMACEval(k2, file_ciphertext)
			if err != nil {
				return contents, errors.New(strings.ToTitle("7"))
			}

			ds_file_HMAC, ok := userlib.DatastoreGet(file_HMAC)
			if !ok {
				return contents, errors.New(strings.ToTitle("Original File Tampered8"))
			}

			check := userlib.HMACEqual(created_file_HMAC, ds_file_HMAC)
			if !check {
				return contents, errors.New(strings.ToTitle("Original File Tampered9"))
			}

			// Decrypt and append
			var forloop_instance ciphertext_index
			marshal_file := userlib.SymDec(k1, file_ciphertext)
			err = json.Unmarshal(marshal_file, &forloop_instance)
			if err != nil {
				return contents, errors.New(strings.ToTitle("7"))
			}
			appended_content := forloop_instance.Content
			contents = append(contents, appended_content...)
		} else {
			bytes_index, err := json.Marshal(i)
			if err != nil {
				return contents, errors.New(strings.ToTitle("7"))
			}
			file_UUID, err = uuid.FromBytes(userlib.Hash([]byte(string(k1) + string(k2) + string(bytes_index)))[:16])
			if err != nil {
				return contents, errors.New(strings.ToTitle("5"))
			}
			file_HMAC, err = uuid.FromBytes(userlib.Hash([]byte(string(k2) + string(k1) + string(bytes_index)))[:16])
			if err != nil {
				return contents, errors.New(strings.ToTitle("6"))
			}

			file_ciphertext, ok := userlib.DatastoreGet(file_UUID)
			if !ok {
				return contents, errors.New(strings.ToTitle("4"))
			}

			created_file_HMAC, err := userlib.HMACEval(k2, file_ciphertext)
			if err != nil {
				return contents, errors.New(strings.ToTitle("7"))
			}

			ds_file_HMAC, ok := userlib.DatastoreGet(file_HMAC)
			if !ok {
				return contents, errors.New(strings.ToTitle("Original File Tampered8"))
			}

			check := userlib.HMACEqual(created_file_HMAC, ds_file_HMAC)
			if !check {
				return contents, errors.New(strings.ToTitle(string(k2)))
			}
			var forloop_instance ciphertext_index
			marshal_file := userlib.SymDec(k1, file_ciphertext)
			err = json.Unmarshal(marshal_file, &forloop_instance)
			if err != nil {
				return contents, errors.New(strings.ToTitle("invalid"))
			}
			appended_content := forloop_instance.Content
			contents = append(contents, appended_content...)
		}
	}

	return contents, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	storage_content_uuid, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("2"))
	}

	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username + filename))[:16])
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("2"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	// checking HMAC
	ciphertext, ok := userlib.DatastoreGet(storage_content_uuid)
	if !ok {
		return uuid.Nil, errors.New(strings.ToTitle("Original File Tampered"))
	}

	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, ciphertext)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("3"))
	}

	ds_HMAC_key, ok := userlib.DatastoreGet(storage_HMAC_uuid)
	if !ok {
		return uuid.Nil, errors.New(strings.ToTitle("Original File Tampered"))
	}

	check := userlib.HMACEqual(ds_HMAC_key, storage_HMAC)
	if !check {
		return uuid.Nil, errors.New(strings.ToTitle("Original File Tampered"))
	}

	//Decrypting ciphertext
	var file File
	marshal_file := userlib.SymDec(filestruct_ekey, ciphertext)
	err = json.Unmarshal(marshal_file, &file)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("Unmarshall Error"))
	}

	k1 := file.K1
	k2 := file.K2
	UUID := file.File_UUID
	HMAC := file.HMAC_UUID

	invitation_uuid := file.File_UUID
	invitation_HMAC_uuid := file.HMAC_UUID
	sym_enc_key := file.K1
	hmac_key := file.K2
	iv := userlib.RandomBytes(16)

	if file.Received {
		invitation_ciphertext, ok := userlib.DatastoreGet(invitation_uuid)
		if !ok {
			return uuid.Nil, errors.New(strings.ToTitle("1"))
		}

		storage_invitation_HMAC, err := userlib.HMACEval(hmac_key, invitation_ciphertext)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("Invalid"))
		}

		ds_HMAC, ok := userlib.DatastoreGet(invitation_HMAC_uuid)
		if !ok {
			return uuid.Nil, errors.New(strings.ToTitle("2"))
		}

		check := userlib.HMACEqual(ds_HMAC, storage_invitation_HMAC)
		if !check {
			return uuid.Nil, errors.New(strings.ToTitle("3"))
		}

		//Decrypting ciphertext
		var invitation Invitation
		marshal_file := userlib.SymDec(sym_enc_key, invitation_ciphertext)
		err = json.Unmarshal(marshal_file, &invitation)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("Unmarshall Error"))
		}

		if invitation.Reject {
			return uuid.Nil, errors.New(strings.ToTitle("Revoked"))
		}

		//encrypt and sign Marshal file
		invitation_uuid = file.File_UUID
		invitation_HMAC_uuid = file.HMAC_UUID
		sym_enc_key = file.K1
		hmac_key = file.K2

	} else {
		//create invitation for the first time
		var invitation Invitation
		invitation.K1 = k1
		invitation.K2 = k2
		invitation.File_UUID = UUID
		invitation.HMAC_UUID = HMAC

		//uploading invitation
		invitation_uuid = uuid.New()
		invitation_HMAC_uuid = uuid.New()
		sym_enc_key = userlib.RandomBytes(16)
		hmac_key = userlib.RandomBytes(16)

		bytes_invitation, err := json.Marshal(invitation)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("9"))
		}

		invitation_ciphertext := userlib.SymEnc(sym_enc_key, iv, bytes_invitation)
		invitation_HMAC, err := userlib.HMACEval(hmac_key, invitation_ciphertext)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("10"))
		}

		userlib.DatastoreSet(invitation_uuid, invitation_ciphertext)
		userlib.DatastoreSet(invitation_HMAC_uuid, invitation_HMAC)

		//update file struct
		file.Recipient_UUID[recipientUsername] = invitation_uuid
		file.Recipient_HMAC[recipientUsername] = invitation_HMAC_uuid
		file.Recipient_K1[recipientUsername] = sym_enc_key
		file.Recipient_k2[recipientUsername] = hmac_key

		byte_new_file, err := json.Marshal(file)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("11"))
		}

		new_file_ciphertext := userlib.SymEnc(filestruct_ekey, iv, byte_new_file)
		new_file_HMAC, err := userlib.HMACEval(filestruct_hkey, new_file_ciphertext)
		if err != nil {
			return uuid.Nil, errors.New(strings.ToTitle("11"))
		}

		userlib.DatastoreSet(storage_content_uuid, new_file_ciphertext)
		userlib.DatastoreSet(storage_HMAC_uuid, new_file_HMAC)
	}
	//creating invitation message
	invitation_message_uuid := uuid.New()

	bytes_invitation_uuid, err := json.Marshal(invitation_uuid)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("11"))
	}
	bytes_invitation_HMAC, err := json.Marshal(invitation_HMAC_uuid)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("11"))
	}
	bytes_k1, err := json.Marshal(sym_enc_key)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("11"))
	}
	bytes_k2, err := json.Marshal(hmac_key)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("11"))
	}

	recipient_pub_key, ok := userlib.KeystoreGet(recipientUsername)
	if !ok {
		return uuid.Nil, errors.New(strings.ToTitle("Original File Tampered"))
	}

	invitation_uuid_ciphertext, err := userlib.PKEEnc(recipient_pub_key, bytes_invitation_uuid)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("error"))
	}

	invitation_HMAC_ciphertext, err := userlib.PKEEnc(recipient_pub_key, bytes_invitation_HMAC)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("error"))
	}

	invitation_k1_ciphertext, err := userlib.PKEEnc(recipient_pub_key, bytes_k1)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("error"))
	}

	invitation_k2_ciphertext, err := userlib.PKEEnc(recipient_pub_key, bytes_k2)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("error"))
	}
	sender_sign_key := userdata.Sign_key

	uuid_dig_sign, err := userlib.DSSign(sender_sign_key, invitation_uuid_ciphertext)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("132"))
	}
	hmac_dig_sign, err := userlib.DSSign(sender_sign_key, invitation_HMAC_ciphertext)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("13"))
	}
	k1_dig_sign, err := userlib.DSSign(sender_sign_key, invitation_k1_ciphertext)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("13"))
	}
	k2_dig_sign, err := userlib.DSSign(sender_sign_key, invitation_k2_ciphertext)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("13"))
	}

	var invitation_info Invitation_info
	invitation_info.Bytes_UUID = invitation_uuid_ciphertext
	invitation_info.Bytes_HMAC = invitation_HMAC_ciphertext
	invitation_info.Bytes_K1 = invitation_k1_ciphertext
	invitation_info.Bytes_K2 = invitation_k2_ciphertext
	invitation_info.Sign_UUID = uuid_dig_sign
	invitation_info.Sign_HMAC = hmac_dig_sign
	invitation_info.Sign_K1 = k1_dig_sign
	invitation_info.Sign_K2 = k2_dig_sign

	bytes_invitation_info, err := json.Marshal(invitation_info)
	if err != nil {
		return uuid.Nil, errors.New(strings.ToTitle("11"))
	}

	userlib.DatastoreSet(invitation_message_uuid, bytes_invitation_info)

	return invitation_message_uuid, nil

}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	var invitation_info Invitation_info
	content_Datastore, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New(strings.ToTitle("error"))
	}

	check1 := json.Unmarshal(content_Datastore, &invitation_info)
	if check1 != nil {
		return errors.New(strings.ToTitle("error"))
	}
	invitation_uuid_ciphertext := invitation_info.Bytes_UUID
	invitation_HMAC_ciphertext := invitation_info.Bytes_HMAC
	invitation_k1_ciphertext := invitation_info.Bytes_K1
	invitation_k2_ciphertext := invitation_info.Bytes_K2
	uuid_dig_sign := invitation_info.Sign_UUID
	hmac_dig_sign := invitation_info.Sign_HMAC
	k1_dig_sign := invitation_info.Sign_K1
	k2_dig_sign := invitation_info.Sign_K2

	verify_key, ok := userlib.KeystoreGet(senderUsername + "ver_key")
	//verify
	check2 := userlib.DSVerify(verify_key, invitation_uuid_ciphertext, uuid_dig_sign)
	if check2 != nil {
		return errors.New(strings.ToTitle("error"))
	}
	check3 := userlib.DSVerify(verify_key, invitation_HMAC_ciphertext, hmac_dig_sign)
	if check3 != nil {
		return errors.New(strings.ToTitle("error"))
	}
	check4 := userlib.DSVerify(verify_key, invitation_k1_ciphertext, k1_dig_sign)
	if check4 != nil {
		return errors.New(strings.ToTitle("error"))
	}
	check5 := userlib.DSVerify(verify_key, invitation_k2_ciphertext, k2_dig_sign)
	if check5 != nil {
		return errors.New(strings.ToTitle("error"))
	}

	//public decryption
	invitation_uuid := uuid.New()
	invitation_HMAC_uuid := uuid.New()
	sym_enc_key := userlib.RandomBytes(16)
	hmac_key := userlib.RandomBytes(16)

	marshal_file_1, err := userlib.PKEDec(userdata.Priv_key, invitation_uuid_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	err = json.Unmarshal(marshal_file_1, &invitation_uuid)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	marshal_file_2, err := userlib.PKEDec(userdata.Priv_key, invitation_HMAC_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	err = json.Unmarshal(marshal_file_2, &invitation_HMAC_uuid)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	marshal_file_3, err := userlib.PKEDec(userdata.Priv_key, invitation_k1_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	err = json.Unmarshal(marshal_file_3, &sym_enc_key)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	marshal_file_4, err := userlib.PKEDec(userdata.Priv_key, invitation_k2_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	err = json.Unmarshal(marshal_file_4, &hmac_key)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	//check if the invitation has been revoked
	invitation_ciphertext, ok := userlib.DatastoreGet(invitation_uuid)
	if !ok {
		return errors.New(strings.ToTitle("1"))
	}

	storage_invitation_HMAC, err := userlib.HMACEval(hmac_key, invitation_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("Invalid"))
	}

	ds_HMAC, ok := userlib.DatastoreGet(invitation_HMAC_uuid)
	if !ok {
		return errors.New(strings.ToTitle("2"))
	}

	check := userlib.HMACEqual(ds_HMAC, storage_invitation_HMAC)
	if !check {
		return errors.New(strings.ToTitle("3"))
	}

	//Decrypting ciphertext
	var invitation Invitation
	marshal_file := userlib.SymDec(sym_enc_key, invitation_ciphertext)
	err = json.Unmarshal(marshal_file, &invitation)
	if err != nil {
		return errors.New(strings.ToTitle("Unmarshall Error"))
	}

	if invitation.Reject {
		return errors.New(strings.ToTitle("Revoked"))
	}

	//file creation
	var file File
	file.File_UUID = invitation_uuid
	file.HMAC_UUID = invitation_HMAC_uuid
	file.K1 = sym_enc_key
	file.K2 = hmac_key
	file.Received = true

	storage_content_uuid, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	// check if the filename already exists
	check_filename, ok := userlib.DatastoreGet(storage_content_uuid)
	if ok {
		return errors.New(strings.ToTitle(string(check_filename)))
	}

	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username + filename))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	iv := userlib.RandomBytes(16)
	bytes_file, err := json.Marshal(file)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	storage_ciphertext := userlib.SymEnc(filestruct_ekey, iv, bytes_file)
	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, storage_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	userlib.DatastoreSet(storage_content_uuid, storage_ciphertext)
	userlib.DatastoreSet(storage_HMAC_uuid, storage_HMAC)

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	storage_content_uuid, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}
	storage_HMAC_uuid, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username + filename))[:16])
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	filestruct_ekey := userlib.Hash([]byte(filename + userdata.Username + "encryption"))[:16]
	filestruct_hkey := userlib.Hash([]byte(filename + userdata.Username + "HMAC"))[:16]

	// checking HMAC
	ciphertext, ok := userlib.DatastoreGet(storage_content_uuid)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	storage_HMAC, err := userlib.HMACEval(filestruct_hkey, ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	ds_HMAC_key, ok := userlib.DatastoreGet(storage_HMAC_uuid)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	check := userlib.HMACEqual(ds_HMAC_key, storage_HMAC)
	if !check {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	//Decrypting ciphertext
	var file File
	marshal_file := userlib.SymDec(filestruct_ekey, ciphertext)
	err = json.Unmarshal(marshal_file, &file)
	if err != nil {
		return errors.New(strings.ToTitle("Unmarshall Error"))
	}
	if file.Received {
		return errors.New(strings.ToTitle("You are not the owener of the file"))
	}
	if file.Recipient_UUID[recipientUsername] == uuid.Nil {
		return errors.New(strings.ToTitle("file not shared with the recipient"))
	}
	//old file info
	old_uuid := file.File_UUID
	old_hmac := file.HMAC_UUID
	old_k1 := file.K1
	old_k2 := file.K2

	old_ciphertext, ok := userlib.DatastoreGet(old_uuid)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	old_HMAC, err := userlib.HMACEval(old_k2, old_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	old_ds_HMAC_key, ok := userlib.DatastoreGet(old_hmac)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	old_check := userlib.HMACEqual(old_HMAC, old_ds_HMAC_key)
	if !old_check {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}
	var file_info ciphertext_index
	marshal_file_info := userlib.SymDec(old_k1, old_ciphertext)
	err = json.Unmarshal(marshal_file_info, &file_info)
	if err != nil {
		return errors.New(strings.ToTitle("7"))
	}

	//create new keys for the file
	k1 := userlib.RandomBytes(16)
	k2 := userlib.RandomBytes(16)
	file_uuid := uuid.New()
	HMAC_uuid := uuid.New()
	iv := userlib.RandomBytes(16)

	//update content with new file info
	index := file_info.Index

	for i := 0; i <= index; i += 1 {
		if i == 0 {
			file_ciphertext := userlib.SymEnc(k1, iv, marshal_file_info)
			file_HMAC, err := userlib.HMACEval(k2, file_ciphertext)
			if err != nil {
				return errors.New(strings.ToTitle("Invalid"))
			}

			userlib.DatastoreSet(file_uuid, file_ciphertext)
			userlib.DatastoreSet(HMAC_uuid, file_HMAC)

		} else {
			bytes_index, err := json.Marshal(i)

			old_file_UUID, err := uuid.FromBytes(userlib.Hash([]byte(string(old_k1) + string(old_k2) + string(bytes_index)))[:16])
			if err != nil {
				return errors.New(strings.ToTitle("5"))
			}
			old_file_HMAC, err := uuid.FromBytes(userlib.Hash([]byte(string(old_k2) + string(old_k1) + string(bytes_index)))[:16])
			if err != nil {
				return errors.New(strings.ToTitle("6"))
			}

			file_ciphertext, ok := userlib.DatastoreGet(old_file_UUID)
			if !ok {
				return errors.New(strings.ToTitle("invalid"))
			}

			created_file_HMAC, err := userlib.HMACEval(old_k2, file_ciphertext)
			if err != nil {
				return errors.New(strings.ToTitle("7"))
			}

			ds_file_HMAC, ok := userlib.DatastoreGet(old_file_HMAC)
			if !ok {
				return errors.New(strings.ToTitle("Original File Tampered8"))
			}

			check := userlib.HMACEqual(created_file_HMAC, ds_file_HMAC)
			if !check {
				return errors.New(strings.ToTitle(string(k2)))
			}
			var file_info ciphertext_index
			marshal_file := userlib.SymDec(old_k1, file_ciphertext)
			err = json.Unmarshal(marshal_file, &file_info)
			if err != nil {
				return errors.New(strings.ToTitle("7"))
			}
			//encrypt
			new_file_ciphertext := userlib.SymEnc(k1, iv, marshal_file)
			file_HMAC, err := userlib.HMACEval(k2, new_file_ciphertext)
			if err != nil {
				return errors.New(strings.ToTitle("Invalid"))
			}

			new_file_UUID, err := uuid.FromBytes(userlib.Hash([]byte(string(k1) + string(k2) + string(bytes_index)))[:16])
			if err != nil {
				return errors.New(strings.ToTitle("5"))
			}
			new_hmac_UUID, err := uuid.FromBytes(userlib.Hash([]byte(string(k2) + string(k1) + string(bytes_index)))[:16])
			if err != nil {
				return errors.New(strings.ToTitle("5"))
			}

			userlib.DatastoreSet(new_file_UUID, new_file_ciphertext)
			userlib.DatastoreSet(new_hmac_UUID, file_HMAC)
		}
	}
	// Reject recipient invitation
	recipient_uuid := file.Recipient_UUID[recipientUsername]
	recipient_hmac := file.Recipient_HMAC[recipientUsername]
	recipient_k1 := file.Recipient_K1[recipientUsername]
	recipient_k2 := file.Recipient_k2[recipientUsername]

	recipient_ciphertext, ok := userlib.DatastoreGet(recipient_uuid)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	recipient_HMAC, err := userlib.HMACEval(recipient_k2, recipient_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	recipient_ds_HMAC_key, ok := userlib.DatastoreGet(recipient_hmac)
	if !ok {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}

	recipient_check := userlib.HMACEqual(recipient_HMAC, recipient_ds_HMAC_key)
	if !recipient_check {
		return errors.New(strings.ToTitle("Original File Tampered"))
	}
	var recipient_invitation Invitation
	marshal_recipient_info := userlib.SymDec(recipient_k1, recipient_ciphertext)
	err = json.Unmarshal(marshal_recipient_info, &recipient_invitation)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	recipient_invitation.Reject = true

	byte_recipient_filedata, err := json.Marshal(recipient_invitation)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	new_recipient_ciphertext := userlib.SymEnc(recipient_k1, iv, byte_recipient_filedata)
	new_recipient_HMAC, err := userlib.HMACEval(recipient_k2, new_recipient_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	userlib.DatastoreSet(recipient_uuid, new_recipient_ciphertext)
	userlib.DatastoreSet(recipient_hmac, new_recipient_HMAC)

	// delete recipient from the map
	file.K1 = k1
	file.K2 = k2
	file.File_UUID = file_uuid
	file.HMAC_UUID = HMAC_uuid
	delete(file.Recipient_UUID, recipientUsername)
	delete(file.Recipient_HMAC, recipientUsername)
	delete(file.Recipient_K1, recipientUsername)
	delete(file.Recipient_k2, recipientUsername)

	byte_filedata, err := json.Marshal(file)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	new_storage_ciphertext := userlib.SymEnc(filestruct_ekey, iv, byte_filedata)
	new_storage_HMAC, err := userlib.HMACEval(filestruct_hkey, new_storage_ciphertext)
	if err != nil {
		return errors.New(strings.ToTitle("error"))
	}

	userlib.DatastoreSet(storage_content_uuid, new_storage_ciphertext)
	userlib.DatastoreSet(storage_HMAC_uuid, new_storage_HMAC)

	//Traverse through the map and update the values
	for key := range file.Recipient_UUID {
		recipient_UUID := file.Recipient_UUID[key]
		recipient_HMAC := file.Recipient_HMAC[key]
		recipient_k1 := file.Recipient_K1[key]
		recipient_k2 := file.Recipient_k2[key]

		ciphertext, ok := userlib.DatastoreGet(recipient_UUID)
		if !ok {
			return errors.New(strings.ToTitle("Original File Tampered"))
		}

		storage_HMAC, err := userlib.HMACEval(recipient_k2, ciphertext)
		if err != nil {
			return errors.New(strings.ToTitle("Invalid"))
		}

		ds_HMAC_key, ok := userlib.DatastoreGet(recipient_HMAC)
		if !ok {
			return errors.New(strings.ToTitle("Original File Tampered"))
		}

		check := userlib.HMACEqual(ds_HMAC_key, storage_HMAC)
		if !check {
			return errors.New(strings.ToTitle("Original File Tampered"))
		}

		//Decrypting ciphertext
		var old_invitation Invitation
		marshal_file := userlib.SymDec(recipient_k1, ciphertext)
		err = json.Unmarshal(marshal_file, &old_invitation)
		if err != nil {
			return errors.New(strings.ToTitle("7"))
		}

		old_invitation.File_UUID = file_uuid
		old_invitation.HMAC_UUID = HMAC_uuid
		old_invitation.K1 = k1
		old_invitation.K2 = k2

		byte_new_invitation, err := json.Marshal(old_invitation)
		if err != nil {
			return errors.New(strings.ToTitle("error"))
		}

		new_invitation_ciphertext := userlib.SymEnc(recipient_k1, iv, byte_new_invitation)
		new_invitation_HMAC, err := userlib.HMACEval(recipient_k2, new_invitation_ciphertext)
		if err != nil {
			return errors.New(strings.ToTitle("error"))
		}

		userlib.DatastoreSet(recipient_UUID, new_invitation_ciphertext)
		userlib.DatastoreSet(recipient_HMAC, new_invitation_HMAC)
	}
	return nil
}
