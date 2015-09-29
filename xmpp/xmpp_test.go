package xmpp

import (
	"encoding/xml"
	"io"
	"net"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type XmppSuite struct{}

var _ = Suite(&XmppSuite{})

func (s *XmppSuite) TestDiscoReplyVerSimple(c *C) {
	expect := "QgayPKawpkPSDYmwT/WM94uAlu0="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://code.google.com/p/exodus#QgayPKawpkPSDYmwT/WM94uAlu0='>
    <identity category='client' name='Exodus 0.9.1' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
  </query>
  `)
	var dr DiscoveryReply
	c.Assert(xml.Unmarshal(input, &dr), IsNil)
	hash, err := dr.VerificationString()
	c.Assert(err, IsNil)
	c.Assert(hash, Equals, expect)
}

func (s *XmppSuite) TestDiscoReplyVerComplex(c *C) {
	expect := "q07IKJEyjvHSyhy//CH0CxmKi8w="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://psi-im.org#q07IKJEyjvHSyhy//CH0CxmKi8w='>
    <identity xml:lang='en' category='client' name='Psi 0.11' type='pc'/>
    <identity xml:lang='el' category='client' name='Ψ 0.11' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
    <x xmlns='jabber:x:data' type='result'>
      <field var='FORM_TYPE' type='hidden'>
        <value>urn:xmpp:dataforms:softwareinfo</value>
      </field>
      <field var='ip_version'>
        <value>ipv4</value>
        <value>ipv6</value>
      </field>
      <field var='os'>
        <value>Mac</value>
      </field>
      <field var='os_version'>
        <value>10.5.1</value>
      </field>
      <field var='software'>
        <value>Psi</value>
      </field>
      <field var='software_version'>
        <value>0.11</value>
      </field>
    </x>
  </query>
`)
	var dr DiscoveryReply
	c.Assert(xml.Unmarshal(input, &dr), IsNil)
	hash, err := dr.VerificationString()
	c.Assert(err, IsNil)
	c.Assert(hash, Equals, expect)
}

type mockConn struct {
	calledClose int
	net.TCPConn
}

func (c *mockConn) Close() error {
	c.calledClose++
	return nil
}
func (s *XmppSuite) TestConnClose(c *C) {
	mockConfigConn := mockConn{}
	conn := Conn{
		config: &Config{
			Conn: &mockConfigConn,
		},
	}
	c.Assert(conn.Close(), IsNil)
	c.Assert(mockConfigConn.calledClose, Equals, 1)
}

type mockConnIOReaderWriter struct {
	read []byte
	err  error
}

func (in mockConnIOReaderWriter) Read(p []byte) (n int, err error) {
	copy(p, in.read)
	return len(in.read), in.err
}

func (s *XmppSuite) TestConnNextEOF(c *C) {
	mockIn := mockConnIOReaderWriter{err: io.EOF}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{})
	c.Assert(stanza.Value, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *XmppSuite) TestConnNextErr(c *C) {
	mockIn := mockConnIOReaderWriter{
		read: []byte(`
      <field var='os'>
        <value>Mac</value>
      </field>
		`),
	}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{})
	c.Assert(stanza.Value, IsNil)
	c.Assert(err.Error(), Equals, "unexpected XMPP message  <field/>")
}

func (s *XmppSuite) TestConnNextIQSet(c *C) {
	mockIn := mockConnIOReaderWriter{
		read: []byte(`
<iq to='example.com'
    xmlns='jabber:client'
    type='set'
    id='sess_1'>
  <session xmlns='urn:ietf:params:xml:ns:xmpp-session'/>
</iq>
  `),
	}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{Space: NsClient, Local: "iq"})
	iq, ok := stanza.Value.(*ClientIQ)
	c.Assert(ok, Equals, true)
	c.Assert(iq.To, Equals, "example.com")
	c.Assert(iq.Type, Equals, "set")
	c.Assert(err, IsNil)
}

func (s *XmppSuite) TestConnNextIQResult(c *C) {
	mockIn := mockConnIOReaderWriter{
		read: []byte(`
<iq from='example.com'
    xmlns='jabber:client'
    type='result'
    id='sess_1'/>
  `),
	}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{Space: NsClient, Local: "iq"})
	iq, ok := stanza.Value.(*ClientIQ)
	c.Assert(ok, Equals, true)
	c.Assert(iq.From, Equals, "example.com")
	c.Assert(iq.Type, Equals, "result")
	c.Assert(err, ErrorMatches, "xmpp: failed to parse id from iq: .*")
}
